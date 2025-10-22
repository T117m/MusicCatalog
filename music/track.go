package music

import (
	"slices"
	"strings"
)

type Track struct {
	ID       int    `db:"id"`
	Title    string `db:"title"`
	Artist   string `db:"artist"`
	Genre    string `db:"genre"`
	FileType string `db:"file_type"`
	FilePath string `db:"file_path"`
}

func New(title, artist, genre, fileType, filePath string) Track {
	return Track{
		Title:    title,
		Artist:   artist,
		Genre:    genre,
		FileType: fileType,
		FilePath: filePath,
	}
}

func (t Track) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return ErrEmptyTitle
	}

	if strings.TrimSpace(t.Artist) == "" {
		return ErrEmptyArtist
	}

	if !t.IsSupportedFormat() {
		return ErrUnsupportedFormat
	}

	if strings.TrimSpace(t.FilePath) == "" {
		return ErrEmptyFilePath
	}

	return nil
}

