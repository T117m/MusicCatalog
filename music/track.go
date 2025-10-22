package music

import (
	"errors"
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

const (
	MP3  = "mp3"
	FLAC = "flac"
	WAV  = "wav"
	M4A  = "m4a"
)

var SupportedFormats = []string{MP3, FLAC, WAV, M4A}

func (t Track) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return errors.New("Title cannot be empty")
	}

	if strings.TrimSpace(t.Artist) == "" {
		return errors.New("Artist cannot be empty")
	}

	if !t.IsSupportedFormat() {
		return errors.New("File type not supported")
	}

	if strings.TrimSpace(t.FilePath) == "" {
		return errors.New("File path cannot be empty")
	}

	return nil
}

func (t Track) IsSupportedFormat() bool {
	return slices.Contains(SupportedFormats, t.FileType)
}
