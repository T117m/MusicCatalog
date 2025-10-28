package storage

import (
	"os"
	"path/filepath"
	"fmt"
	"github.com/T117m/MusicCatalog/music"
)

func (s *Storage) AddTrack(track music.Track) error {
	if err := track.Validate(); err != nil {
		return fmt.Errorf("couldn't add track: %w", err)
	}

	q := "INSERT INTO tracks(title, artist, genre, file_type, file_path) VALUES (?, ?, ?, ?, ?);"

	s.db.Exec(q, track.Title, track.Artist, track.Genre, track.Genre, track.FileType, track.FilePath)

	return nil
}

func (s *Storage) GetAllTracks() ([]music.Track, error) {
	return nil, nil
}

func (s *Storage) GetTracksByArtist(artist string) ([]music.Track, error) {
	return nil, nil
}
