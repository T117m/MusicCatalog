package music

import (
	"path/filepath"
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

	if strings.TrimSpace(t.Genre) == "" {
		return ErrEmptyGenre
	}

	if t.FileType == "" {
		return ErrEmptyFileType
	}

	if strings.TrimSpace(t.FilePath) == "" {
		return ErrEmptyFilePath
	}

	if !t.IsSupportedFormat() {
		return ErrUnsupportedFormat
	}

	return nil
}

func (t *Track) Normalize() {
	if strings.TrimSpace(t.Title) == "" {
		t.Title = "Untitled"
	}

	if strings.TrimSpace(t.Artist) == "" {
		t.Artist = "Unknown"
	}

	if strings.TrimSpace(t.Genre) == "" {
		t.Genre = "Other"
	}

	if t.FileType == "" {
		if ext := strings.TrimPrefix(filepath.Ext(t.FilePath), "."); ext != "" {
			t.FileType = strings.ToLower(ext)
		}
	}
}
