package storage

import (
	"path/filepath"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	dbPath := filepath.Join("storage", "internal", "catalog.db")

	db, err :=  sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	storage := &Storage{db: db}

	if err := storage.db.Migrate(); err != nil {
		fmt.Errorf("failed to migrate database: %w", err)
	}

	return storage, nil
}
