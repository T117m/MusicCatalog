package ui

import (
	"log"
	"strconv"
	"strings"

	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/key"
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
	EditTrackView
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
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, len(m.inputs))
	)
	m.tracks, cmd = m.tracks.Update(msg)
	cmds = append(cmds, cmd)
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		
		if pressed := msg.String(); pressed == "ctrl+c" {
			return m, tea.Quit
		}

		switch m.view {
		case TrackListView:
			switch {
			case key.Matches(msg, trackListKeyMap.next):
				m.tracks.MoveDown(1)
			case key.Matches(msg, trackListKeyMap.prev):
				m.tracks.MoveUp(1)
			case key.Matches(msg, trackListKeyMap.quit):
				return m, tea.Quit
			case key.Matches(msg, trackListKeyMap.add):
				m.tracks.Blur()
				m.view = AddTrackView
				cmd := m.inputs[m.focused].Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, trackListKeyMap.edit):
				m.tracks.Blur()
				m.view = EditTrackView
				cmd := m.inputs[m.focused].Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, trackListKeyMap.delete):
				m.tracks.Blur()
				m.view = DeleteTrackView
			case msg.String() == "enter":
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
						//m.player.Wait()
						m.player.Play(&track)
					}
				}
			}
		case AddTrackView:
			switch {
			case key.Matches(msg, addTrackKeyMap.next):
				m.nextInput()
			case key.Matches(msg, addTrackKeyMap.prev):
				m.inputs[m.focused].Blur()
				m.prevInput()
				cmd = m.inputs[m.focused].Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, addTrackKeyMap.quit):
				m.quitInput()
			case key.Matches(msg, addTrackKeyMap.add):
				m.addTrack()
			}
		case EditTrackView:
			switch {
			case key.Matches(msg, editTrackKeyMap.next):
				m.nextInput()
			case key.Matches(msg, editTrackKeyMap.prev):
				m.inputs[m.focused].Blur()
				m.prevInput()
				cmd = m.inputs[m.focused].Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, editTrackKeyMap.quit):
				m.quitInput()
			case key.Matches(msg, editTrackKeyMap.edit):
				m.editTrack()
			}
		case DeleteTrackView:
			switch {
			case key.Matches(msg, deleteTrackKeyMap.quit):
				m.view = TrackListView
				m.tracks.Focus()
			case key.Matches(msg, deleteTrackKeyMap.delete):
				if m.errMsg == nil {
				m.removeTrack()
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
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
		sb.WriteString(gloss.PlaceHorizontal(
			gloss.Width(baseStyle.Render(m.tracks.View()))+
				gloss.Width(helpStyle.Render(inputHelp)),
			gloss.Right,
			helpStyle.Render(inputHelp),
		))
	case DeleteTrackView:
		sb.WriteString(m.renderDeletePrompt())
		sb.WriteString(helpStyle.Render(deleteHelp))
	case EditTrackView:
		sb.WriteString(
			gloss.JoinHorizontal(
				gloss.Top,
				baseStyle.Render(m.tracks.View()),
				baseStyle.Width(30).Render(m.renderInputForm()),
			))
		sb.WriteString(gloss.PlaceHorizontal(
			gloss.Width(baseStyle.Render(m.tracks.View()))+
				gloss.Width(helpStyle.Render(inputHelp)),
			gloss.Right,
			helpStyle.Render(inputHelp),
		))
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
		m.errMsg = nil
		m.tracks = newTrackList(m.storage)

		m.quitInput()
	}
}

func (m *model) removeTrack() {
	id, _ := strconv.Atoi(m.tracks.SelectedRow()[0])

	if err := m.storage.RemoveTrackByID(id); err != nil {
		m.errMsg = err
	} else {
		m.errMsg = nil
		m.tracks = newTrackList(m.storage)
		m.view = TrackListView

		m.tracks.Focus()
	}
}

func (m *model) editTrack() {
	id, _ := strconv.Atoi(m.tracks.SelectedRow()[0])
	newTitle, newArtist, newGenre, newFT, newFP := m.getInputs()

	if err := m.storage.EditTrackByID(id, newTitle, newArtist, newGenre, newFT, newFP); err != nil {
		m.errMsg = err

		if err == music.ErrEmptyFilePath || err == music.ErrUnsupportedFormat {
			m.setFocus(fp)
		} else {
			log.Fatal(err)
		}
	} else {
		m.errMsg = nil
		m.tracks = newTrackList(m.storage)

		m.quitInput()
	}
}
