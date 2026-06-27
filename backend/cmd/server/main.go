package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/make-smart-products/speed-match/internal/config"
	"github.com/make-smart-products/speed-match/internal/db"
	"github.com/make-smart-products/speed-match/internal/handler"
	"github.com/make-smart-products/speed-match/internal/middleware"
	"github.com/make-smart-products/speed-match/internal/repository"
	"github.com/make-smart-products/speed-match/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		log.Fatalf("db dir: %v", err)
	}
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		log.Fatalf("upload dir: %v", err)
	}

	conn, err := db.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer conn.Close()

	if err := runMigrations(conn.DB); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := repository.New(conn)
	events := service.NewEventService(repo)
	h := handler.New(events, cfg.UploadDir)

	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.CORS(cfg.CORSOrigin))
	r.Use(middleware.JSONContentType)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))
	r.Mount("/", h.Routes(middleware.RateLimit(10)))

	staticDir := resolveStaticDir(cfg.StaticDir)
	if staticDir != "" {
		log.Printf("serving frontend from %s", staticDir)
		spa := spaHandler(staticDir)
		r.NotFound(spa.ServeHTTP)
		r.MethodNotAllowed(spa.ServeHTTP)
	} else {
		log.Printf("STATIC_DIR not found (%q), API-only mode", cfg.StaticDir)
	}

	addr := ":" + cfg.Port
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func runMigrations(sqlDB *sql.DB) error {
	migrationsDir := findMigrationsDir()
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}
	return goose.Up(sqlDB, migrationsDir)
}

func findMigrationsDir() string {
	candidates := []string{
		"migrations",
		"backend/migrations",
		filepath.Join("..", "migrations"),
	}
	for _, c := range candidates {
		if dirExists(c) {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return "migrations"
}

func resolveStaticDir(path string) string {
	if path == "" {
		return ""
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return ""
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func spaHandler(staticDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(staticDir))
	indexPath := filepath.Join(staticDir, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/uploads/") {
			http.NotFound(w, r)
			return
		}

		if r.URL.Path != "/" {
			candidate := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
			if rel, err := filepath.Rel(staticDir, candidate); err == nil && !strings.HasPrefix(rel, "..") {
				if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
					fileServer.ServeHTTP(w, r)
					return
				}
			}
		}

		http.ServeFile(w, r, indexPath)
	})
}
