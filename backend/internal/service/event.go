package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/make-smart-products/speed-match/internal/models"
	"github.com/make-smart-products/speed-match/internal/repository"
)

var (
	ErrVotingClosed   = errors.New("voting is not open")
	ErrResultsClosed  = errors.New("results are not available yet")
	ErrVoteLimit      = errors.New("vote limit exceeded")
	ErrSelfVote       = errors.New("cannot vote for yourself")
	ErrInvalidTargets = errors.New("invalid target participants")
)

type EventService struct {
	repo *repository.Repository
}

func NewEventService(repo *repository.Repository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(title string, voteLimit *int) (*models.Event, string, error) {
	baseSlug := Slugify(title)
	slug := UniqueSlug(baseSlug, func(sl string) bool {
		ok, err := s.repo.SlugExists(sl)
		return err == nil && ok
	})

	adminToken, err := NewToken()
	if err != nil {
		return nil, "", err
	}

	event := &models.Event{
		ID:         NewID(),
		Title:      title,
		Slug:       slug,
		VoteLimit:  voteLimit,
		Status:     models.StatusDraft,
		AdminToken: adminToken,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}

	if err := s.repo.CreateEvent(event); err != nil {
		return nil, "", err
	}

	return event, adminToken, nil
}

func (s *EventService) AddParticipant(eventID, pseudonym string, photoURL *string) (*models.Participant, error) {
	token, err := NewToken()
	if err != nil {
		return nil, err
	}

	p := &models.Participant{
		ID:          NewID(),
		EventID:     eventID,
		Pseudonym:   pseudonym,
		PhotoURL:    photoURL,
		AccessToken: token,
	}

	if err := s.repo.CreateParticipant(p); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *EventService) GetStatus(slug, accessToken string) (*models.EventStatusResponse, error) {
	event, err := s.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, err
	}

	participant, err := s.repo.GetParticipantByToken(accessToken)
	if err != nil {
		return nil, err
	}
	if participant.EventID != event.ID {
		return nil, repository.ErrNotFound
	}

	selected, err := s.repo.GetVotesByVoter(participant.ID)
	if err != nil {
		return nil, err
	}

	return &models.EventStatusResponse{
		Title:     event.Title,
		Slug:      event.Slug,
		Status:    event.Status,
		VoteLimit: event.VoteLimit,
		Selected:  selected,
		Pseudonym: participant.Pseudonym,
	}, nil
}

func (s *EventService) ListParticipants(slug, accessToken string) ([]models.Participant, error) {
	event, err := s.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, err
	}

	participant, err := s.repo.GetParticipantByToken(accessToken)
	if err != nil {
		return nil, err
	}
	if participant.EventID != event.ID {
		return nil, repository.ErrNotFound
	}

	return s.repo.ListOtherParticipants(event.ID, participant.ID)
}

func (s *EventService) SaveVotesByToken(accessToken string, targetIDs []string) error {
	participant, err := s.repo.GetParticipantByToken(accessToken)
	if err != nil {
		return err
	}

	event, err := s.repo.GetEventByID(participant.EventID)
	if err != nil {
		return err
	}

	return s.saveVotes(event, participant, targetIDs)
}

func (s *EventService) SaveVotes(slug, accessToken string, targetIDs []string) error {
	event, err := s.repo.GetEventBySlug(slug)
	if err != nil {
		return err
	}

	participant, err := s.repo.GetParticipantByToken(accessToken)
	if err != nil {
		return err
	}
	if participant.EventID != event.ID {
		return repository.ErrNotFound
	}

	return s.saveVotes(event, participant, targetIDs)
}

func (s *EventService) saveVotes(event *models.Event, participant *models.Participant, targetIDs []string) error {
	if event.Status != models.StatusVoting {
		return ErrVotingClosed
	}

	for _, id := range targetIDs {
		if id == participant.ID {
			return ErrSelfVote
		}
	}

	if event.VoteLimit != nil && len(targetIDs) > *event.VoteLimit {
		return fmt.Errorf("%w: max %d", ErrVoteLimit, *event.VoteLimit)
	}

	if len(targetIDs) > 0 {
		count, err := s.repo.CountParticipantsInEvent(event.ID, targetIDs)
		if err != nil {
			return err
		}
		if count != len(targetIDs) {
			return ErrInvalidTargets
		}
	}

	return s.repo.ReplaceVotes(event.ID, participant.ID, targetIDs)
}

func (s *EventService) ListMatches(slug, accessToken string) ([]models.Participant, error) {
	event, err := s.repo.GetEventBySlug(slug)
	if err != nil {
		return nil, err
	}

	if event.Status != models.StatusClosed {
		return nil, ErrResultsClosed
	}

	participant, err := s.repo.GetParticipantByToken(accessToken)
	if err != nil {
		return nil, err
	}
	if participant.EventID != event.ID {
		return nil, repository.ErrNotFound
	}

	return s.repo.ListMatches(participant.ID, event.ID)
}

func (s *EventService) GetAdminEvent(adminToken, slug string) (*models.AdminEventResponse, error) {
	event, err := s.repo.GetEventByAdminToken(adminToken)
	if err != nil {
		return nil, err
	}
	if event.Slug != slug {
		return nil, repository.ErrNotFound
	}

	participants, err := s.repo.ListParticipants(event.ID)
	if err != nil {
		return nil, err
	}

	return &models.AdminEventResponse{
		Event:        *event,
		Participants: participants,
	}, nil
}

func (s *EventService) PatchEvent(adminToken, slug string, req models.PatchEventRequest) (*models.Event, error) {
	event, err := s.repo.GetEventByAdminToken(adminToken)
	if err != nil {
		return nil, err
	}
	if event.Slug != slug {
		return nil, repository.ErrNotFound
	}

	if req.VoteLimit != nil {
		event.VoteLimit = req.VoteLimit
	}
	if req.Status != nil {
		event.Status = *req.Status
	}

	if err := s.repo.UpdateEvent(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) GetEventByAdminToken(adminToken string) (*models.Event, error) {
	return s.repo.GetEventByAdminToken(adminToken)
}
