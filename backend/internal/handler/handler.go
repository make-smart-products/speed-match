package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/make-smart-products/speed-match/internal/models"
	"github.com/make-smart-products/speed-match/internal/repository"
	"github.com/make-smart-products/speed-match/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const maxUploadSize = 2 << 20 // 2 MB

type Handler struct {
	events    *service.EventService
	uploadDir string
}

func New(events *service.EventService, uploadDir string) *Handler {
	return &Handler{events: events, uploadDir: uploadDir}
}

func (h *Handler) RegisterRoutes(r chi.Router, voteRateLimit func(http.Handler) http.Handler) {
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/events/{slug}/status", h.GetStatus)
		api.Get("/events/{slug}/participants", h.GetParticipants)
		api.Get("/events/{slug}/matches", h.GetMatches)
		api.With(voteRateLimit).Post("/votes", h.SaveVotes)

		api.Post("/admin/events", h.CreateEvent)
		api.Get("/admin/events/{slug}", h.GetAdminEvent)
		api.Post("/admin/events/{slug}/participants", h.AddParticipant)
		api.Patch("/admin/events/{slug}", h.PatchEvent)
	})
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing access token")
		return
	}

	resp, err := h.events.GetStatus(slug, token)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetParticipants(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing access token")
		return
	}

	list, err := h.events.ListParticipants(slug, token)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) GetMatches(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing access token")
		return
	}

	list, err := h.events.ListMatches(slug, token)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	if list == nil {
		list = []models.Participant{}
	}

	writeJSON(w, http.StatusOK, list)
}

func (h *Handler) SaveVotes(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Access-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing access token")
		return
	}

	var req models.SaveVotesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.TargetIDs == nil {
		req.TargetIDs = []string{}
	}

	if err := h.events.SaveVotesByToken(token, req.TargetIDs); err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	event, adminToken, err := h.events.CreateEvent(req.Title, req.VoteLimit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create event")
		return
	}

	writeJSON(w, http.StatusCreated, models.CreateEventResponse{
		Event:      *event,
		AdminToken: adminToken,
		AdminURL:   fmt.Sprintf("/admin/%s?key=%s", event.Slug, adminToken),
	})
}

func (h *Handler) GetAdminEvent(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	token := r.Header.Get("X-Admin-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing admin token")
		return
	}

	resp, err := h.events.GetAdminEvent(token, slug)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) AddParticipant(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	token := r.Header.Get("X-Admin-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing admin token")
		return
	}

	event, err := h.events.GetEventByAdminToken(token)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	if event.Slug != slug {
		writeError(w, http.StatusNotFound, "event not found")
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "invalid form data")
		return
	}

	pseudonym := strings.TrimSpace(r.FormValue("pseudonym"))
	if pseudonym == "" {
		writeError(w, http.StatusBadRequest, "pseudonym is required")
		return
	}

	var photoURL *string
	file, header, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		data, err := io.ReadAll(io.LimitReader(file, maxUploadSize+1))
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to read photo")
			return
		}
		if len(data) > maxUploadSize {
			writeError(w, http.StatusBadRequest, "photo too large (max 2MB)")
			return
		}

		ext, err := detectImageExt(data, header.Filename)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		eventDir := filepath.Join(h.uploadDir, event.ID)
		if err := os.MkdirAll(eventDir, 0o755); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save photo")
			return
		}

		filename := uuid.NewString() + ext
		fullPath := filepath.Join(eventDir, filename)
		if err := os.WriteFile(fullPath, data, 0o644); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save photo")
			return
		}

		url := "/uploads/" + event.ID + "/" + filename
		photoURL = &url
	}

	p, err := h.events.AddParticipant(event.ID, pseudonym, photoURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to add participant")
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) PatchEvent(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	token := r.Header.Get("X-Admin-Token")
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing admin token")
		return
	}

	var req models.PatchEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	event, err := h.events.PatchEvent(token, slug, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, event)
}

func detectImageExt(data []byte, filename string) (string, error) {
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg", nil
	}
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return ".png", nil
	}
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") {
		return ".jpg", nil
	}
	if strings.HasSuffix(lower, ".png") {
		return ".png", nil
	}
	return "", errors.New("only jpeg and png images are allowed")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, service.ErrVotingClosed):
		writeError(w, http.StatusConflict, "voting is not open")
	case errors.Is(err, service.ErrResultsClosed):
		writeError(w, http.StatusConflict, "results are not available yet")
	case errors.Is(err, service.ErrVoteLimit):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrSelfVote):
		writeError(w, http.StatusBadRequest, "cannot vote for yourself")
	case errors.Is(err, service.ErrInvalidTargets):
		writeError(w, http.StatusBadRequest, "invalid target participants")
	default:
		log.Printf("internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal error")
	}
}
