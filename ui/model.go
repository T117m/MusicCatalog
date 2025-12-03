package ui

import (
	"strconv"
	"strings"

	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type model struct {
	storage *storage.Storage
	player  *player.Player
	tracks  table.Model
	view    ViewMode
	inputs  []ti.Model
	focused int
	errMsg  error
}

type ViewMode int

const (
	TrackListView ViewMode = iota
	AddTrackView
	DeleteTrackView
	PlayerView
)

func New(store *storage.Storage, player *player.Player) model {
	tracks := newTrackList(store)
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
		cmds = make([]tea.Cmd, len(m.inputs))
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.view == TrackListView {
				return m, tea.Quit
			}
		case "esc":
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
				cmd := m.inputs[m.focused].Focus()
				cmds = append(cmds, cmd)
			case TrackListView:
				m.tracks.MoveUp(1)
			}
		case "ctrl+a":
			if m.view == TrackListView {
				m.tracks.Blur()
				m.view = AddTrackView
				cmd := m.inputs[m.focused].Focus()
				cmds = append(cmds, cmd)
			}
		case "ctrl+s":
			if m.view == AddTrackView {
				m.addTrack()
			}
		case "x":
			if m.view == TrackListView {
				m.tracks.Blur()
				m.view = DeleteTrackView
			}
		}
	}

	switch m.view {
	case TrackListView:
		var cmd tea.Cmd
		m.tracks, cmd = m.tracks.Update(msg)

		return m, cmd
	case AddTrackView:
		for i := range m.inputs {
			var cmd tea.Cmd
			m.inputs[i], cmd = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	switch m.view {
	case TrackListView:
		sb.WriteString(baseStyle.Render(m.tracks.View()))
		sb.WriteString(helpStyle.Render(trackListHelp))
	case AddTrackView:
		sb.WriteString(
			gloss.JoinHorizontal(
				gloss.Top,
				baseStyle.Render(m.tracks.View()),
				baseStyle.Width(30).Render(m.renderInputForm()),
			))
		sb.WriteString(helpStyle.Render(inputHelp))
	case DeleteTrackView:
		sb.WriteString(
			gloss.JoinHorizontal(
				gloss.Top,
				baseStyle.Render(m.tracks.View()),
				m.renderDeletePrompt(),
			))
	case PlayerView:
	}

	return sb.String()
}

func (m *model) addTrack() {
	newTrack := music.New(m.getInputs())

	newTrack.Normalize()

	if err := newTrack.Validate(); err != nil {
		if err == music.ErrEmptyFilePath || err == music.ErrUnsupportedFormat {
			m.errMsg = err
			m.setFocus(fp)
		}
	} else if err := m.storage.AddTrack(&newTrack); err != nil {
			m.errMsg = err
			m.setFocus(fp)
	} else {
		m.tracks = newTrackList(m.storage)
		
		m.quitInput()
	}
}
