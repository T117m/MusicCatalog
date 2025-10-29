package main

import (
	"fmt"
	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/storage"
	"log"
)

func main() {
	track := music.New("INTERLOPER (seventeen_split.cfg) #JerseyClub", "DJ CARCINOGEN", "", "mp3", "INTERLOPER (seventeen_split.cfg) #JerseyClub - DJ CARCINOGEN.mp3")
	track.Normalize()

	if err := track.Validate(); err != nil {
		log.Fatalf("error creating track: %s", err)
	}

	fmt.Printf("трек создан: %+v\n", track)

	strg, err := storage.New()
	if err != nil {
		log.Fatalf("error creating storage: %s", err)
	}

	if err = strg.AddTrack(track); err != nil {
		log.Fatalf("error adding track: %s", err)
	}
}
