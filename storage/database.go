package storage

import (
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

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening a database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	storage := &Storage{db: db}

	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("error reading migrations directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		migration, err := migrations.ReadFile(filepath.Join("migrations", name))
		if err != nil {
			return nil, fmt.Errorf("error reading migration %s: %w", name, err)
		}

		_, err = storage.db.Exec(string(migration))
		if err != nil {
			return nil, fmt.Errorf("error executing migration %s: %w", name, err)
		}
	}

	return storage, nil
}
