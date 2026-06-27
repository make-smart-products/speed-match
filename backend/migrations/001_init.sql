-- +goose Up
CREATE TABLE events (
    id             TEXT PRIMARY KEY,
    title          TEXT NOT NULL,
    slug           TEXT NOT NULL UNIQUE,
    vote_limit     INTEGER,
    status         TEXT NOT NULL DEFAULT 'draft',
    admin_token    TEXT NOT NULL UNIQUE,
    created_at     TEXT NOT NULL DEFAULT (datetime('now')),
    voting_ends_at TEXT
);

CREATE TABLE participants (
    id           TEXT PRIMARY KEY,
    event_id     TEXT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    pseudonym    TEXT NOT NULL,
    photo_url    TEXT,
    access_token TEXT NOT NULL UNIQUE,
    UNIQUE (event_id, pseudonym)
);

CREATE TABLE votes (
    id         TEXT PRIMARY KEY,
    event_id   TEXT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    voter_id   TEXT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    target_id  TEXT NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE (voter_id, target_id)
);

CREATE INDEX idx_participants_event ON participants(event_id);
CREATE INDEX idx_votes_voter ON votes(voter_id);
CREATE INDEX idx_votes_event ON votes(event_id);

-- +goose Down
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS events;
