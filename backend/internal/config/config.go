package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBPath     string
	UploadDir  string
	CORSOrigin string
	StaticDir  string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/speedmatch.db"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		if os.Getenv("STATIC_DIR") != "" || isProduction() {
			corsOrigin = "*"
		} else {
			corsOrigin = "http://localhost:5173"
		}
	}

	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "../web/dist"
	}

	return Config{
		Port:       port,
		DBPath:     dbPath,
		UploadDir:  uploadDir,
		CORSOrigin: corsOrigin,
		StaticDir:  staticDir,
	}
}

func isProduction() bool {
	return os.Getenv("RENDER") != "" || os.Getenv("RAILWAY_ENVIRONMENT") != "" || os.Getenv("ENV") == "production"
}

func GetEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}
