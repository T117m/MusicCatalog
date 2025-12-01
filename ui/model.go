package ui

import (
	"log"
	"strconv"
	"strings"

	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	storage *storage.Storage
	player  *player.Player
	tracks  table.Model
	view    ViewMode
	inputs  []ti.Model
	focused int
}

type ViewMode int

const (
	TrackListView ViewMode = iota
	AddTrackView
	PlayerView
)

func New(store *storage.Storage, player *player.Player) model {
	tracks := newTracksTable(store)
	inputs := newInputs()

	return model{
		storage: store,
		player:  player,
		tracks:  tracks,
		view:    TrackListView,
		inputs:  inputs,
		focused: 0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Музыкальный католог")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, len(m.inputs))
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			switch m.view {
			case TrackListView:
				return m, tea.Quit
			case AddTrackView:
				m.quitInput()
			}
		case "enter":
			switch m.view {
			case TrackListView:
				if len(m.tracks.Rows()) > 0 {
					id, _ := strconv.Atoi(m.tracks.SelectedRow()[0])

					track, _ := m.storage.GetTrackByID(id)

					switch m.player.GetState() {
					case player.Playing:
						if m.player.GetCurrentTrack().ID == id {
							m.player.Pause()
						} else {
							m.player.Stop()
							m.player.Wait()
							m.player.Play(&track)
						}
					case player.Paused:
						if m.player.GetCurrentTrack().ID == id {
							m.player.Play(&track)
						} else {
							m.player.Stop()
							m.player.Wait()
							m.player.Play(&track)
						}
					default:
						m.player.Play(&track)
					}
				}
			case AddTrackView:
				if m.focused == len(m.inputs)-1 {
					m.addTrack()
				} else {
					m.nextInput()
				}
			}
		case "tab", "ctrl+n":
			switch m.view {
			case AddTrackView:
				m.nextInput()
			case TrackListView:
				m.tracks.MoveDown(1)
			}
		case "shift+tab", "ctrl+p":
			switch m.view {
			case AddTrackView:
				m.inputs[m.focused].Blur()
				if m.view == AddTrackView {
					m.prevInput()
				}
				m.inputs[m.focused].Focus()
			case TrackListView:
				m.tracks.MoveUp(1)
			}
		case "a", "i":
			if m.view == TrackListView {
				m.tracks.Blur()
				m.view = AddTrackView
				m.resetInputs()
		}
		case "ctrl+s":
			if m.view == AddTrackView {
				m.addTrack()
			}
		}
	}

	switch m.view {
	case TrackListView:
		m.tracks, cmd = m.tracks.Update(msg)

		return m, cmd
	case AddTrackView:
		for i := range m.inputs {
			m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() string {
	var s strings.Builder

	switch m.view {
	case TrackListView:
		s.WriteString(baseStyle.Render(m.tracks.View()))
		s.WriteString(helpStyle.Render(trackListHelp))
	case AddTrackView:
		s.WriteString(baseStyle.Render(m.tracks.View()))
		s.WriteString("\n")
		s.WriteString(baseStyle.Render(renderInputForm(m.inputs)))
		s.WriteString(helpStyle.Render(inputHelp))
	case PlayerView:
	}

	return s.String()
}

func newTracksTable(store *storage.Storage) table.Model {
	tracks, _ := store.GetAllTracks()

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Название", Width: 12},
		{Title: "Исполнитель", Width: 12},
		{Title: "Тип файла", Width: 10},
		{Title: "Жанр", Width: 10},
	}

	var rows []table.Row
	for _, track := range tracks {
		row := []string{
			strconv.Itoa(track.ID),
			track.Title,
			track.Artist,
			track.FileType,
			track.Genre,
		}

		rows = append(rows, row)
	}

	return newStyledTable(columns, rows)
}

func (m *model) addTrack() {
	newTrack := music.New(m.getInputs())

	newTrack.Normalize()

	if err := newTrack.Validate(); err != nil {
		log.Fatal(err)
		m.quitInput()
	}

	m.storage.AddTrack(&newTrack)
	m.quitInput()
}
