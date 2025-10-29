package storage

import (
	"os"
	"fmt"
	"github.com/T117m/MusicCatalog/music"
)

func (s *Storage) AddTrack(track music.Track) error {
	if _, err := os.Stat(track.FilePath); err != nil {
		return fmt.Errorf("can't add track: %w", err)
	}

	track.Normalize()

	if err := track.Validate(); err != nil {
		return fmt.Errorf("can't add track: validation fail: %w", err)
	}

	q := "INSERT INTO tracks(title, artist, genre, file_type, file_path) VALUES (?, ?, ?, ?, ?);"

	_, err := s.db.Exec(q, track.Title, track.Artist, track.Genre, track.FileType, track.FilePath)
	if err != nil {
		return fmt.Errorf("can't add track due to query error: %w", err)
	}

	return nil
}

func (s *Storage) GetAllTracks() ([]music.Track, error) {
	return nil, nil
}

func (s *Storage) GetTracksByArtist(artist string) ([]music.Track, error) {
	return nil, nil
}
