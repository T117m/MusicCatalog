package music

import (
	"slices"
	"strings"
)

const (
	MP3  = "mp3"
	FLAC = "flac"
	WAV  = "wav"
	M4A  = "m4a"
)

var SupportedFormats = []string{MP3, FLAC, WAV, M4A}

func (t Track) IsSupportedFormat() bool {
	f := strings.ToLower(t.FileType)
	return slices.Contains(SupportedFormats, f)
}
