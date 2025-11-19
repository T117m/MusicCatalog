package player

import (
	"fmt"
	"github.com/T117m/MusicCatalog/music"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"os"
	"time"
)

type Player struct {
	state        PlayerState
	currentTrack *music.Track
	source       beep.StreamSeekCloser
	streamer     beep.Streamer
	format       beep.Format
	ctrl         *beep.Ctrl
	done         chan bool
	stop         chan struct{}
}

type PlayerState int

const (
	Stopped PlayerState = iota
	Playing
	Paused
)

const defaultSampleRate = beep.SampleRate(44100)

func New() *Player {
	speaker.Init(defaultSampleRate, defaultSampleRate.N(time.Second/10))

	return &Player{
		state: Stopped,
		done:  make(chan bool),
		stop:  make(chan struct{}),
	}
}

func (p *Player) Play(track *music.Track) error {
	f, err := os.Open(track.FilePath)
	if err != nil {
		return fmt.Errorf("can't open %s: %w", track.FilePath, err)
	}

	var (
		streamer beep.StreamSeekCloser
		format   beep.Format
	)

	switch track.FileType {
	case music.MP3:
		streamer, format, err = mp3.Decode(f)
	default:
		return music.ErrUnsupportedFormat
	}

	if err != nil {
		p.source.Close()
		return fmt.Errorf("can't decode %s: %w", track.FilePath, err)
	}

	p.source = streamer

	if format.SampleRate != defaultSampleRate {
		p.streamer = beep.Resample(4, format.SampleRate, defaultSampleRate, streamer)
	} else {
		p.streamer = streamer
	}

	p.format = format
	p.currentTrack = track
	p.ctrl = &beep.Ctrl{Streamer: streamer}
	p.state = Playing

	go func() {
		speaker.Play(beep.Seq(p.ctrl, beep.Callback(func() {
			p.done <- true
		})))

		select {
		case <-p.done:
			p.state = Stopped
		case <-p.stop:
			speaker.Lock()
			p.ctrl.Streamer = nil
			speaker.Unlock()
		}

		if p.source != nil {
			p.source.Close()
		}
	}()

	return nil
}

func (p *Player) Wait() {
	<-p.done
}

func (p *Player) IsPlaying() bool {
	return p.state == Playing
}

func (p *Player) GetCurrentTrack() *music.Track {
	return p.currentTrack
}

