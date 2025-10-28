package storage

import (
	"os"
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	dbPath := filepath.Join("storage", "internal", "catalog.db")

	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening a database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	storage := &Storage{db: db}

	migration, err := migrations.ReadFile("migrations/001_create_tracks.sql")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error reading migration: %w", err)
	}

	_, err = storage.db.Exec(string(migration))
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error executing migration: %w", err)
	}

	return storage, nil
}
