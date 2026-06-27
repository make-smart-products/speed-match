package models

type EventStatus string

const (
	StatusDraft  EventStatus = "draft"
	StatusVoting EventStatus = "voting"
	StatusClosed EventStatus = "closed"
)

type Event struct {
	ID           string      `db:"id" json:"id"`
	Title        string      `db:"title" json:"title"`
	Slug         string      `db:"slug" json:"slug"`
	VoteLimit    *int        `db:"vote_limit" json:"vote_limit"`
	Status       EventStatus `db:"status" json:"status"`
	AdminToken   string      `db:"admin_token" json:"-"`
	CreatedAt    string      `db:"created_at" json:"created_at"`
	VotingEndsAt *string     `db:"voting_ends_at" json:"voting_ends_at,omitempty"`
}

type Participant struct {
	ID          string  `db:"id" json:"id"`
	EventID     string  `db:"event_id" json:"event_id"`
	Pseudonym   string  `db:"pseudonym" json:"pseudonym"`
	PhotoURL    *string `db:"photo_url" json:"photo_url,omitempty"`
	AccessToken string  `db:"access_token" json:"access_token,omitempty"`
}

type Vote struct {
	ID       string `db:"id"`
	EventID  string `db:"event_id"`
	VoterID  string `db:"voter_id"`
	TargetID string `db:"target_id"`
}

type EventStatusResponse struct {
	Title      string      `json:"title"`
	Slug       string      `json:"slug"`
	Status     EventStatus `json:"status"`
	VoteLimit  *int        `json:"vote_limit"`
	Selected   []string    `json:"selected_ids"`
	Pseudonym  string      `json:"pseudonym"`
}

type CreateEventRequest struct {
	Title     string `json:"title"`
	VoteLimit *int   `json:"vote_limit"`
}

type CreateEventResponse struct {
	Event      Event  `json:"event"`
	AdminToken string `json:"admin_token"`
	AdminURL   string `json:"admin_url"`
}

type SaveVotesRequest struct {
	TargetIDs []string `json:"target_ids"`
}

type PatchEventRequest struct {
	Status    *EventStatus `json:"status"`
	VoteLimit *int         `json:"vote_limit"`
}

type AdminEventResponse struct {
	Event        Event         `json:"event"`
	Participants []Participant `json:"participants"`
}
