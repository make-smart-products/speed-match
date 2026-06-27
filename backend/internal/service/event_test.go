package service_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/make-smart-products/speed-match/internal/db"
	"github.com/make-smart-products/speed-match/internal/models"
	"github.com/make-smart-products/speed-match/internal/repository"
	"github.com/make-smart-products/speed-match/internal/service"
	"github.com/pressly/goose/v3"
)

func TestMutualMatchFlow(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	conn, err := db.Connect(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	if err := goose.SetDialect("sqlite3"); err != nil {
		t.Fatal(err)
	}
	migrationsDir, _ := filepath.Abs("../../migrations")
	if err := goose.Up(conn.DB, migrationsDir); err != nil {
		t.Fatal(err)
	}

	repo := repository.New(conn)
	svc := service.NewEventService(repo)

	event, _, err := svc.CreateEvent("Test", intPtr(3))
	if err != nil {
		t.Fatal(err)
	}

	p1, err := svc.AddParticipant(event.ID, "Alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := svc.AddParticipant(event.ID, "Bob", nil)
	if err != nil {
		t.Fatal(err)
	}

	voting := models.StatusVoting
	if _, err := svc.PatchEvent(event.AdminToken, event.Slug, models.PatchEventRequest{Status: &voting}); err != nil {
		t.Fatal(err)
	}

	if err := svc.SaveVotesByToken(p1.AccessToken, []string{p2.ID}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SaveVotesByToken(p2.AccessToken, []string{p1.ID}); err != nil {
		t.Fatal(err)
	}

	closed := models.StatusClosed
	if _, err := svc.PatchEvent(event.AdminToken, event.Slug, models.PatchEventRequest{Status: &closed}); err != nil {
		t.Fatal(err)
	}

	matches, err := svc.ListMatches(event.Slug, p1.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 1 || matches[0].Pseudonym != "Bob" {
		t.Fatalf("expected Bob match, got %+v", matches)
	}

	matches2, err := svc.ListMatches(event.Slug, p2.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches2) != 1 || matches2[0].Pseudonym != "Alice" {
		t.Fatalf("expected Alice match, got %+v", matches2)
	}

	// one-sided vote should not appear
	p3, err := svc.AddParticipant(event.ID, "Carol", nil)
	if err != nil {
		t.Fatal(err)
	}
	open := models.StatusVoting
	if _, err := svc.PatchEvent(event.AdminToken, event.Slug, models.PatchEventRequest{Status: &open}); err != nil {
		t.Fatal(err)
	}
	if err := svc.SaveVotesByToken(p3.AccessToken, []string{p1.ID}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.PatchEvent(event.AdminToken, event.Slug, models.PatchEventRequest{Status: &closed}); err != nil {
		t.Fatal(err)
	}
	matches3, err := svc.ListMatches(event.Slug, p1.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if len(matches3) != 1 {
		t.Fatalf("one-sided sympathy must not create extra matches, got %d", len(matches3))
	}
}

func intPtr(v int) *int { return &v }

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
