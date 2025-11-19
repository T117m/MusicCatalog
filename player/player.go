package player

import (
	"fmt"
	"github.com/T117m/MusicCatalog/music"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
	"os"
	"sync"
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
	mutex        sync.Mutex
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
		done:  make(chan bool, 1),
		stop:  make(chan struct{}),
	}
}

func (p *Player) Play(track *music.Track) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.state == Playing {
		p.stopPlayback()
	}

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
	p.ctrl = &beep.Ctrl{Streamer: p.streamer}
	p.state = Playing

	p.done = make(chan bool, 1)
	p.stop = make(chan struct{})

	go func() {
		speaker.Play(beep.Seq(p.ctrl, beep.Callback(func() {
			select {
			case p.done <- true:
			default:
			}
		})))

		select {
		case <-p.done:
			p.mutex.Lock()
			if p.state == Playing {
				p.state = Stopped
			}
			p.mutex.Unlock()
		case <-p.stop:
			speaker.Lock()
			if p.ctrl != nil {
				p.ctrl.Streamer = nil
			}
			speaker.Unlock()
		}

		p.mutex.Lock()
		if p.source != nil {
			p.source.Close()
			p.source = nil
		}
		p.mutex.Unlock()
	}()

	return nil
}

func (p *Player) stopPlayback() {
	if p.stop != nil {
		close(p.stop)
		p.stop = nil
	}
	p.state = Stopped
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
