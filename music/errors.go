package music

import "errors"

var (
	ErrEmptyTitle        = errors.New("title cannot be empty")
	ErrEmptyArtist       = errors.New("artist cannot be empty")
	ErrEmptyFileType     = errors.New("file type cannot be empty")
	ErrUnsupportedFormat = errors.New("format not supported")
	ErrEmptyFilePath     = errors.New("file path cannot be empty")
	ErrEmptyGenre        = errors.New("genre cannot be empty")
)
