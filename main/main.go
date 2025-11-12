package main

import (
	"fmt"
	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/storage"
	"log"
)

func main() {
	track := music.New("test", "", "", "mp3", "storage/internal/test.mp3")
	track.Normalize()

	if err := track.Validate(); err != nil {
		log.Fatalf("error creating track: %s", err)
	}

	fmt.Printf("трек создан: %+v\n", track)

	strg, err := storage.New()
	if err != nil {
		log.Fatalf("error creating storage: %s", err)
	}

	if err = strg.AddTrack(&track); err != nil {
		log.Fatalf("error adding track: %s", err)
	}

	tracks, err := strg.GetAllTracks()
	if err != nil {
		log.Fatalf("error getting all tracks: %s", err)
	}

	fmt.Printf("треки извлечены: %+v\n", tracks)

	unknowns, err := strg.GetTracksByArtist("Unknown")
	if err != nil {
		log.Fatalf("error getting tracks by Unknown: %s", err)
	}

	fmt.Printf("треки от Unknown: %+v\n", unknowns)

	knowns, err := strg.GetTracksByArtist("Known")
	if err != nil {
		log.Fatalf("error getting tracks by Known: %s", err)
	}

	fmt.Printf("треки от Known: %+v\n", knowns)

	idTrack, err := strg.GetTrackByID(track.ID)
	if err != nil {
		log.Fatalf("error getting track by ID %d: %s", track.ID, err)
	}

	fmt.Printf("трек с ID %d: %+v\n", track.ID, idTrack)

	// idTrackF, err := strg.GetTrackByID(0)
	// if err != nil {
		// log.Fatalf("error getting track by ID %d: %s", 0, err)
	// }
// 
	// fmt.Printf("трек с ID %d: %+v\n", 0, idTrackF)

	err = strg.RemoveTrackByID(0)
	if err != nil {
		log.Fatalf("error deleting track by ID %d: %s", 0, err)
	}

	_, err = strg.GetTrackByID(0)
	if err != nil {
		log.Fatalf("error getting track by ID %d: %s", 0, err)
	}
}
