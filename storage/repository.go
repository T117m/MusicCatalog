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
	q := "SELECT * FROM tracks;"

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("can't get tracks: %w", err)
	}
	defer rows.Close()

	tracks := make([]music.Track, 0)
	for rows.Next() {
		track := music.Track{}

		err := rows.Scan(&track.ID, &track.Title, &track.Artist, &track.Genre, &track.FileType, &track.FilePath)
		if err != nil {
			return nil, fmt.Errorf("can't get track: %w", err)
		}

		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't iterate over rows: %w", err)
	}

	return tracks, nil
}

func (s *Storage) GetTracksByArtist(artist string) ([]music.Track, error) {
	return nil, nil
}
