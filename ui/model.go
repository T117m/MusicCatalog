package ui

import (
	"log"
	"strconv"
	"strings"

	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/player"
	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	storage *storage.Storage
	player  *player.Player
	tracks  table.Model
	view    ViewMode
	input   inputFormModel
	search  searchModel
	errMsg  error
	help    help.Model
}

type ViewMode int

const (
	TrackListView ViewMode = iota
	AddTrackView
	DeleteTrackView
	EditTrackView
	SearchView
)

func New(store *storage.Storage, player *player.Player) model {
	var (
		tracks = newTrackList(store)
		input  = newInputModel()
		search = newSearchModel()
	)

	return model{
		storage: store,
		player:  player,
		tracks:  tracks,
		view:    TrackListView,
		input:   input,
		search:  search,
		help:    help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("Музыкальный католог")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds = make([]tea.Cmd, len(m.input.inputs))
	)

	m.tracks, cmd = m.tracks.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	m.search, cmd = m.search.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch {
		case key.Matches(msg, ctrlc):
			return m, tea.Quit
		case key.Matches(msg, showHelp):
			m.help.ShowAll = !m.help.ShowAll
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
			case key.Matches(msg, trackListKeyMap.action):
				if len(m.tracks.Rows()) > 0 {
					id, err := strconv.Atoi(m.tracks.SelectedRow()[0])
					if err != nil {
						m.errMsg = err
					}

					track, err := m.storage.GetTrackByID(id)
					if err != nil {
						m.errMsg = err
					}

					switch m.player.GetState() {
					case player.Playing:
						if m.player.GetCurrentTrack().ID == id {
							m.player.Pause()
						} else {
							m.player.Stop()
							m.player.Wait()
							err = m.player.Play(&track)
							if err != nil {
								m.errMsg = err
							}
						}
					case player.Paused:
						if m.player.GetCurrentTrack().ID == id {
							err = m.player.Play(&track)
							if err != nil {
								m.errMsg = err
							}
						} else {
							m.player.Stop()
							m.player.Wait()
							err = m.player.Play(&track)
							if err != nil {
								m.errMsg = err
							}
						}
					default:
						//m.player.Wait()
						m.player.Play(&track)
					}
				}
			case key.Matches(msg, trackListSubKeyMap.add):
				m.tracks.Blur()
				m.view = AddTrackView
				cmd := m.input.Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, trackListSubKeyMap.edit):
				m.tracks.Blur()
				m.view = EditTrackView
				cmd := m.input.Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, trackListSubKeyMap.delete):
				m.tracks.Blur()
				m.view = DeleteTrackView
			case key.Matches(msg, trackListSubKeyMap.search):
				m.tracks.Blur()
				m.view = SearchView
				m.search.input.Focus()
			}
		case AddTrackView:
			switch {
			case key.Matches(msg, addTrackKeyMap.next):
				m.input.nextInput()
			case key.Matches(msg, addTrackKeyMap.prev):
				m.input.Blur()
				m.input.prevInput()
				cmd = m.input.Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, addTrackKeyMap.quit):
				m.quitInput()
			case key.Matches(msg, addTrackKeyMap.action):
				m.addTrack()
			}
		case EditTrackView:
			switch {
			case key.Matches(msg, editTrackKeyMap.next):
				m.input.nextInput()
			case key.Matches(msg, editTrackKeyMap.prev):
				m.input.Blur()
				m.input.prevInput()
				cmd = m.input.Focus()
				cmds = append(cmds, cmd)
			case key.Matches(msg, editTrackKeyMap.quit):
				m.quitInput()
			case key.Matches(msg, editTrackKeyMap.action):
				m.editTrack()
			}
		case DeleteTrackView:
			switch {
			case key.Matches(msg, deleteTrackKeyMap.quit):
				m.view = TrackListView
				m.tracks.Focus()
			case key.Matches(msg, deleteTrackKeyMap.action):
				if m.errMsg == nil {
					m.removeTrack()
				}
			}
		case SearchView:
			switch {
			case key.Matches(msg, searchTrackKeyMap.next):
				m.search.tag++
				if m.search.tag > 3 {
					m.search.tag = 0
				}
			case key.Matches(msg, searchTrackKeyMap.prev):
				m.search.tag--
				if m.search.tag < 0 {
					m.search.tag = 3
				}
			case key.Matches(msg, searchTrackKeyMap.quit):
				m.search.input.Blur()
				m.view = TrackListView
				m.tracks.Focus()
			case key.Matches(msg, searchTrackKeyMap.action):
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var sb strings.Builder

	switch m.view {
	case TrackListView:
		sb.WriteString(m.search.View(unfocusedStyle))
		sb.WriteString("\n" + baseStyle.Render(m.tracks.View()))
		if m.help.ShowAll {
			sb.WriteString("\n" + m.help.View(trackListSubKeyMap) + "\n")
		}
		sb.WriteString("\n" + m.help.View(trackListKeyMap))
	case AddTrackView:
		m.writeInputForm(&sb, addTrackKeyMap)
	case DeleteTrackView:
		sb.WriteString(m.renderDeletePrompt())
		sb.WriteString("\n" + m.help.View(deleteTrackKeyMap))
	case EditTrackView:
		m.writeInputForm(&sb, editTrackKeyMap)
	case SearchView:
		sb.WriteString(m.search.View(baseStyle))
		sb.WriteString("\n" + unfocusedStyle.Render(m.tracks.View()))
		sb.WriteString("\n" + m.help.View(searchTrackKeyMap))
	}

	return sb.String()
}

func (m *model) addTrack() {
	newTrack := music.New(m.input.getInputs())

	newTrack.Normalize()

	if err := newTrack.Validate(); err != nil {
		if err == music.ErrEmptyFilePath || err == music.ErrUnsupportedFormat {
			m.errMsg = err
			m.input.setFocus(fp)
		}
	} else if err := m.storage.AddTrack(&newTrack); err != nil {
		m.errMsg = err
		m.input.setFocus(fp)
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
	newTitle, newArtist, newGenre, newFT, newFP := m.input.getInputs()

	if err := m.storage.EditTrackByID(id, newTitle, newArtist, newGenre, newFT, newFP); err != nil {
		m.errMsg = err

		if err == music.ErrEmptyFilePath || err == music.ErrUnsupportedFormat {
			m.input.setFocus(fp)
		} else {
			log.Fatal(err)
		}
	} else {
		m.errMsg = nil
		m.tracks = newTrackList(m.storage)

		m.quitInput()
	}
}
