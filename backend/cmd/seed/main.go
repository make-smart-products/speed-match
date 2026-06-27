package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/make-smart-products/speed-match/internal/config"
	"github.com/make-smart-products/speed-match/internal/db"
	"github.com/make-smart-products/speed-match/internal/models"
	"github.com/make-smart-products/speed-match/internal/repository"
	"github.com/make-smart-products/speed-match/internal/service"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	conn, err := db.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer conn.Close()

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "backend/migrations"
	}
	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}
	if err := goose.Up(conn.DB, migrationsDir); err != nil {
		log.Fatal(err)
	}

	repo := repository.New(conn)
	events := service.NewEventService(repo)

	voteLimit := 3
	event, adminToken, err := events.CreateEvent("Демо: Быстрые знакомства", &voteLimit)
	if err != nil {
		log.Fatalf("create event: %v", err)
	}

	names := []string{"Алиса", "Борис", "Вика", "Глеб", "Дина", "Егор"}
	var tokens []string
	for _, name := range names {
		p, err := events.AddParticipant(event.ID, name, nil)
		if err != nil {
			log.Fatalf("add participant %s: %v", name, err)
		}
		tokens = append(tokens, p.AccessToken)
	}

	voting := models.StatusVoting
	event, err = events.PatchEvent(adminToken, event.Slug, models.PatchEventRequest{Status: &voting})
	if err != nil {
		log.Fatalf("open voting: %v", err)
	}

	_ = time.Now()

	fmt.Println("=== Speed Match Demo Seed ===")
	fmt.Printf("Event: %s (%s)\n", event.Title, event.Slug)
	fmt.Printf("Admin URL: /admin/%s?key=%s\n", event.Slug, adminToken)
	fmt.Println()
	fmt.Println("Participant links:")
	for i, name := range names {
		fmt.Printf("  %s: /e/%s?t=%s\n", name, event.Slug, tokens[i])
	}
}
