package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/make-smart-products/speed-match/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SlugExists(slug string) (bool, error) {
	var n int
	err := r.db.Get(&n, `SELECT COUNT(1) FROM events WHERE slug = ?`, slug)
	return n > 0, err
}

func (r *Repository) CreateEvent(e *models.Event) error {
	_, err := r.db.NamedExec(`
		INSERT INTO events (id, title, slug, vote_limit, status, admin_token, created_at)
		VALUES (:id, :title, :slug, :vote_limit, :status, :admin_token, :created_at)
	`, e)
	return err
}

func (r *Repository) GetEventByID(id string) (*models.Event, error) {
	var e models.Event
	err := r.db.Get(&e, `SELECT * FROM events WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &e, err
}

func (r *Repository) GetEventBySlug(slug string) (*models.Event, error) {
	var e models.Event
	err := r.db.Get(&e, `SELECT * FROM events WHERE slug = ?`, slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &e, err
}

func (r *Repository) GetEventByAdminToken(token string) (*models.Event, error) {
	var e models.Event
	err := r.db.Get(&e, `SELECT * FROM events WHERE admin_token = ?`, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &e, err
}

func (r *Repository) UpdateEvent(e *models.Event) error {
	_, err := r.db.NamedExec(`
		UPDATE events SET vote_limit = :vote_limit, status = :status, voting_ends_at = :voting_ends_at
		WHERE id = :id
	`, e)
	return err
}

func (r *Repository) GetParticipantByToken(token string) (*models.Participant, error) {
	var p models.Participant
	err := r.db.Get(&p, `SELECT * FROM participants WHERE access_token = ?`, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &p, err
}

func (r *Repository) CreateParticipant(p *models.Participant) error {
	_, err := r.db.NamedExec(`
		INSERT INTO participants (id, event_id, pseudonym, photo_url, access_token)
		VALUES (:id, :event_id, :pseudonym, :photo_url, :access_token)
	`, p)
	return err
}

func emptySlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func (r *Repository) ListParticipants(eventID string) ([]models.Participant, error) {
	var list []models.Participant
	err := r.db.Select(&list, `
		SELECT id, event_id, pseudonym, photo_url, access_token
		FROM participants WHERE event_id = ? ORDER BY pseudonym
	`, eventID)
	return emptySlice(list), err
}

func (r *Repository) ListOtherParticipants(eventID, selfID string) ([]models.Participant, error) {
	var list []models.Participant
	err := r.db.Select(&list, `
		SELECT id, event_id, pseudonym, photo_url
		FROM participants WHERE event_id = ? AND id != ? ORDER BY pseudonym
	`, eventID, selfID)
	return emptySlice(list), err
}

func (r *Repository) GetVotesByVoter(voterID string) ([]string, error) {
	var ids []string
	err := r.db.Select(&ids, `SELECT target_id FROM votes WHERE voter_id = ?`, voterID)
	return emptySlice(ids), err
}

func (r *Repository) ReplaceVotes(eventID, voterID string, targetIDs []string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM votes WHERE voter_id = ?`, voterID); err != nil {
		return err
	}

	for _, targetID := range targetIDs {
		_, err := tx.Exec(`
			INSERT INTO votes (id, event_id, voter_id, target_id, created_at)
			VALUES (?, ?, ?, ?, ?)
		`, uuid.NewString(), eventID, voterID, targetID, time.Now().UTC().Format(time.RFC3339))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) ListMatches(voterID, eventID string) ([]models.Participant, error) {
	var list []models.Participant
	err := r.db.Select(&list, `
		SELECT p.id, p.event_id, p.pseudonym, p.photo_url
		FROM votes v1
		JOIN votes v2 ON v1.voter_id = v2.target_id AND v1.target_id = v2.voter_id
		JOIN participants p ON p.id = v1.target_id
		WHERE v1.voter_id = ? AND v1.event_id = ?
		ORDER BY p.pseudonym
	`, voterID, eventID)
	return emptySlice(list), err
}

func (r *Repository) CountParticipantsInEvent(eventID string, ids []string) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	query, args, err := sqlx.In(`
		SELECT COUNT(1) FROM participants WHERE event_id = ? AND id IN (?)
	`, eventID, ids)
	if err != nil {
		return 0, err
	}
	query = r.db.Rebind(query)
	var n int
	err = r.db.Get(&n, query, args...)
	return n, err
}

