package music

import (
	"slices"
	"strings"
)

const (
	MP3  = "mp3"
	FLAC = "flac"
	WAV  = "wav"
	OGG  = "ogg"
)

var SupportedFormats = []string{MP3, FLAC, WAV, OGG}

func (t Track) IsSupportedFormat() bool {
	f := strings.ToLower(t.FileType)
	return slices.Contains(SupportedFormats, f)
}
