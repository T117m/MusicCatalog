package storage

import "github.com/T117m/MusicCatalog/music"

func (s *Storage) AddTrack(track music.Track) error {
	return nil
}

func (s *Storage) GetAllTracks() ([]music.Track, error) {
	return nil, nil
}

func (s *Storage) GetTracksByArtist(artist string) ([]music.Track, error) {
	return nil, nil
}
