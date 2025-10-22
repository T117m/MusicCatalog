package music

import "slices"

const (
	MP3  = "mp3"
	FLAC = "flac"
	WAV  = "wav"
	M4A  = "m4a"
)

var SupportedFormats = []string{MP3, FLAC, WAV, M4A}

func (t Track) IsSupportedFormat() bool {
	return slices.Contains(SupportedFormats, t.FileType)
}
