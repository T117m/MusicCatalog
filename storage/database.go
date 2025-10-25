package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"
	"io/fs"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	dbPath := filepath.Join("storage", "internal", "catalog.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening a database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	storage := &Storage{db: db}

	// TODO: make an actual migrations pasing
	q, err := migrations.ReadFile("001_create_tracks.sql")
	if err != nil {
		return nil, fmt.Errorf("error reading migration 001: %w", err)
	}

	_, err = storage.db.Exec(string(q))
	if err != nil {
		return nil, fmt.Errorf("error executing migrations: %w", err)
	}

	return storage, nil
}
