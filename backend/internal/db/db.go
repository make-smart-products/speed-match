package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func Connect(dbPath string) (*sqlx.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", dbPath)
	conn, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	conn.SetMaxOpenConns(1)
	return conn, nil
}
