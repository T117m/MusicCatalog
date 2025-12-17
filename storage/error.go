package storage

import (
	"errors"
	"strings"

	// sqlite "github.com/mattn/go-sqlite3"
)

var (
	ErrFileNotExists = errors.New("file does not exist")
	ErrNoFileAccess  = errors.New("can't access file")
	ErrIsDirectory   = errors.New("is a directory")
	ErrNotUnique     = errors.New("UNIQUE constraint failed")
)

func isUniqueConstraintErr(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
