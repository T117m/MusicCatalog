package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/T117m/MusicCatalog/music"
	"io/fs"
	"os"
	"slices"
	"strings"
)

const (
	insertTrackQuery       = "INSERT INTO tracks(title, artist, genre, file_type, file_path) VALUES (?, ?, ?, ?, ?) RETURNING id;"
	selectAllQuery         = "SELECT id, title, artist, genre, file_type, file_path FROM tracks;"
	selectAllByArtistQuery = "SELECT id, title, artist, genre, file_type, file_path FROM tracks WHERE artist=?;"
	selectByIDQuery        = "SELECT id, title, artist, genre, file_type, file_path FROM tracks WHERE id=?;"
	deleteByIDQuery        = "DELETE FROM tracks WHERE id=?;"
	updateByIDQuery        = "UPDATE tracks SET title=?, artist=?, genre=?, file_type=?, file_path=? WHERE id=?;"
)

func (s *Storage) AddTrack(track *music.Track) error {
	err := checkFilePath(track.FilePath)
	if err != nil {
		return fmt.Errorf("problem with path: %w", err)
	}

	track.Normalize()

	if err := track.Validate(); err != nil {
		return fmt.Errorf("can't add track: validation fail: %w", err)
	}

	err = s.db.QueryRow(insertTrackQuery, track.Title, track.Artist, track.Genre, track.FileType, track.FilePath).Scan(&track.ID)
	if err != nil {
		return fmt.Errorf("can't add track due to query error: %w", err)
	}

	return nil
}

func checkFilePath(fp string) error {
	fileInfo, err := os.Stat(fp)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("file %s does not exist: %w", fp, err)
		}
		return fmt.Errorf("can't access file: %w", err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("%s is a directory: %w", fp, err)
	}

	return nil
}

func (s *Storage) scanTracks(rows *sql.Rows) ([]music.Track, error) {
	defer rows.Close()

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
	rows, err := s.db.Query(selectAllQuery)
	if err != nil {
		return nil, fmt.Errorf("can't get tracks: %w", err)
	}

	return s.scanTracks(rows)
}

func (s *Storage) GetTracksByArtist(artist string) ([]music.Track, error) {
	rows, err := s.db.Query(selectAllByArtistQuery, artist)
	if err != nil {
		return nil, fmt.Errorf("can't get tracks: %w", err)
	}

	return s.scanTracks(rows)
}

func (s *Storage) GetTrackByID(id int) (music.Track, error) {
	track := music.Track{}

	row := s.db.QueryRow(selectByIDQuery, id)
	err := row.Scan(&track.ID, &track.Title, &track.Artist, &track.Genre, &track.FileType, &track.FilePath)
	if err != nil {
		return track, fmt.Errorf("can't scan track: %w", err)
	}

	return track, nil
}

func (s *Storage) RemoveTrackByID(id int) error {
	r, err := s.db.Exec(deleteByIDQuery, id)
	if err != nil {
		return fmt.Errorf("can't delete track by id %d: %w", id, err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("track with id %d not found", id)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) EditTrackByID(id int, title, artist, genre, ft, fp string) error {
	track, err := s.GetTrackByID(id)
	if err != nil {
		return err
	}

	if strings.TrimSpace(title) != "" {
		track.Title = title
	}
	if strings.TrimSpace(artist) != "" {
		track.Artist = artist
	}
	if strings.TrimSpace(genre) != "" {
		track.Genre = genre
	}
	if strings.TrimSpace(ft) != "" && slices.Contains(music.SupportedFormats, ft) {
		track.FileType = ft
	}
	if strings.TrimSpace(fp) != "" {
		if err := checkFilePath(fp); err != nil {
			return err
		}

		track.FilePath = fp
	}

	track.Normalize()
	if err := track.Validate(); err != nil {
		return fmt.Errorf("validation failed after edit: %w", err)
	}

	r, err := s.db.Exec(updateByIDQuery, track.Title, track.Artist, track.Genre, track.FileType, track.FilePath, id)
	if err != nil {
		return fmt.Errorf("can't update track by id %d: %w", id, err)
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("can't get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("track with id %d not found", id)
	}

	return nil
}
