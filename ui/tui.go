package ui

import (
	"fmt"

	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	storage  *storage.Storage
	player   *player.Player
	tracks   []music.Track
	selected int
	view     ViewMode
}

type ViewMode int

const (
	TrackListView ViewMode = iota
	PlayerView
)

func New(store *storage.Storage, player *player.Player) model {
	tracks, _ := store.GetAllTracks()
	return model{
		storage: store,
		player:  player,
		tracks:  tracks,
		view:    TrackListView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			m.selected++
		case "enter", " ":
			if len(m.tracks) > 0 {
				track := m.tracks[m.selected]
				if err := m.player.Play(&track); err != nil {
					// TODO: error handling
				}
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if len(m.tracks) == 0 {
		return "Загрузка треков..."
	}

	var s string

	for i, track := range m.tracks {
		cursor := " "
		if i == m.selected {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s - %s\n", cursor, track.Artist, track.Title)
	}

	s += "\nНажмите q чтобы выйти"

	return s
}
