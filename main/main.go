package main

import (
	"github.com/T117m/MusicCatalog/music"
	// "github.com/T117m/MusicCatalog/storage"
	"github.com/T117m/MusicCatalog/player"
	"log"
)

func main() {
	track := music.New("test", "", "", "mp3", "storage/internal/test.mp3")
	track.Normalize()

	if err := track.Validate(); err != nil {
		log.Fatalf("error creating track: %s", err)
	}

	log.Printf("трек создан: %+v\n", track)

	// strg, err := storage.New()
	// if err != nil {
		// log.Fatalf("error creating storage: %s", err)
	// }
// 
	// if err = strg.AddTrack(&track); err != nil {
		// log.Fatalf("error adding track: %s", err)
	// } else {
		// log.Printf("трек %s добавлен", track.Title)
	// }
// 
	// tracks, err := strg.GetAllTracks()
	// if err != nil {
		// log.Fatalf("error getting all tracks: %s", err)
	// }
// 
	// log.Printf("треки: %+v\n", tracks)

	plr := player.New()

	err := plr.Play(&track)
	if err != nil {
		log.Fatalf("error playing track: %s", err)
	}

	plr.Wait()
}
