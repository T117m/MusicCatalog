package music

import "errors"

var (
	ErrEmptyTitle        = errors.New("Title cannot be empty")
	ErrEmptyArtist       = errors.New("Artist cannot be empty")
	ErrUnsupportedFormat = errors.New("Format not supported")
	ErrEmptyFilePath     = errors.New("File path cannot be empty")
)
