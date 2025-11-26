package ui

import (
	"strconv"

	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	storage  *storage.Storage
	player   *player.Player
	tracks   table.Model
	view     ViewMode
}

type ViewMode int

const (
	TrackListView ViewMode = iota
	PlayerView
)

func New(store *storage.Storage, player *player.Player) model {
	tracks, _ := store.GetAllTracks()

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Artist", Width: 10},
		{Title: "Title", Width: 10},
		{Title: "FileType", Width: 8},
		{Title: "Genre", Width: 10},
	}

	var rows []table.Row
	for _, track := range tracks {
		row := []string{
			strconv.Itoa(track.ID),
			track.Artist,
			track.Title,
			track.FileType, 
			track.Genre,
		}

		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	return model{
		storage: store,
		player:  player,
		tracks:  t,
		view:    TrackListView,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			if len(m.tracks.Rows()) > 0 {
				id, _ := strconv.Atoi(m.tracks.SelectedRow()[0])

				track, _ := m.storage.GetTrackByID(id)

				if err := m.player.Play(&track); err != nil {
					// TODO: error handling
				}
			}
		}
	}

	m.tracks, cmd = m.tracks.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return m.tracks.View() + "\n"
}
