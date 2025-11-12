package storage

import (
	"database/sql"
	"fmt"
	"github.com/T117m/MusicCatalog/music"
	"os"
)

func (s *Storage) AddTrack(track *music.Track) error {
	if _, err := os.Stat(track.FilePath); err != nil {
		return fmt.Errorf("can't add track: %w", err)
	}

	track.Normalize()

	if err := track.Validate(); err != nil {
		return fmt.Errorf("can't add track: validation fail: %w", err)
	}

	q := "INSERT INTO tracks(title, artist, genre, file_type, file_path) VALUES (?, ?, ?, ?, ?) RETURNING id;"

	err := s.db.QueryRow(q, track.Title, track.Artist, track.Genre, track.FileType, track.FilePath).Scan(&track.ID)
	if err != nil {
		return fmt.Errorf("can't add track due to query error: %w", err)
	}

	return nil
}

func (s *Storage) scanTracks(rows *sql.Rows) ([]music.Track, error) {
	tracks := make([]music.Track, 0)
	for rows.Next() {
		track := music.Track{}

		err := rows.Scan(&track.ID, &track.Title, &track.Artist, &track.Genre, &track.FileType, &track.FilePath)
		if err != nil {
			return nil, fmt.Errorf("can't scan track: %w", err)
		}

		tracks = append(tracks, track)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tracks, nil
}

func (s *Storage) GetAllTracks() ([]music.Track, error) {
	q := "SELECT id, title, artist, genre, file_type, file_path FROM tracks;"

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("can't get tracks: %w", err)
	}
	defer rows.Close()

	return s.scanTracks(rows)
}

func (s *Storage) GetTracksByArtist(artist string) ([]music.Track, error) {
	q := "SELECT id, title, artist, genre, file_type, file_path FROM tracks WHERE artist=?;"

	rows, err := s.db.Query(q, artist)
	if err != nil {
		return nil, fmt.Errorf("can't get tracks: %w", err)
	}
	defer rows.Close()

	return s.scanTracks(rows)
}

func (s *Storage) GetTrackByID(id int) (music.Track, error) {
	q := "SELECT id, title, artist, genre, file_type, file_path FROM tracks WHERE id=?;"

	track := music.Track{}

	row := s.db.QueryRow(q, id)
	err := row.Scan(&track.ID, &track.Title, &track.Artist, &track.Genre, &track.FileType, &track.FilePath)
	if err != nil {
		return track, fmt.Errorf("can't scan track: %w", err)
	}

	return track, nil
}

func (s *Storage) RemoveTrackByID(id int) error {
	if _, err := s.GetTrackByID(id); err != nil {
		return fmt.Errorf("can't find track with id %d: %w", id, err)
	}

	q := "DELETE FROM tracks WHERE id=?;"

	_, err := s.db.Exec(q, id)
	if err != nil {
		return fmt.Errorf("can't delete track by id %d: %w", id, err)
	}

	return nil
}
