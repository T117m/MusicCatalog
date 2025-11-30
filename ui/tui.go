package ui

import (
	"fmt"
	"log"
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
}

type ViewMode int

const (
	TrackListView ViewMode = iota
	AddTrackView
	PlayerView
)

var (
	baseStyle = gloss.NewStyle().BorderStyle(gloss.ThickBorder()).
			BorderForeground(gloss.Color("7"))
	helpStyle        = gloss.NewStyle().Foreground(gloss.Color("241"))
	inputHeaderStyle = gloss.NewStyle().Foreground(gloss.Color("7")).Bold(true)
	inputStyle       = gloss.NewStyle().BorderStyle(gloss.NormalBorder()).BorderBottom(true).
				BorderForeground(gloss.Color("7"))
)

const (
	artist = iota
	title
	genre
	ft
	fp
)

const (
	trackListHelp = "\nq: Выйти | a/i: Добавить трек | d: Удалить трек | r: Редактировать трек\n"
	inputHelp     = "\nq: Вернуться | Enter: Ввод\n"
)

func New(store *storage.Storage, player *player.Player) model {
	t := newTracksTable(store)

	inputs := make([]ti.Model, 5)

	inputs[artist] = ti.New()
	inputs[artist].Focus()
	inputs[artist].CharLimit = 20
	inputs[artist].Width = 20
	inputs[artist].Prompt = ""

	inputs[title] = ti.New()
	inputs[title].CharLimit = 20
	inputs[title].Width = 20
	inputs[title].Prompt = ""

	inputs[genre] = ti.New()
	inputs[genre].CharLimit = 10
	inputs[genre].Width = 10
	inputs[genre].Prompt = ""

	inputs[ft] = ti.New()
	inputs[ft].CharLimit = 4
	inputs[ft].Width = 5
	inputs[ft].Prompt = ""

	inputs[fp] = ti.New()
	inputs[fp].CharLimit = 30
	inputs[fp].Width = 30
	inputs[fp].Prompt = ""

	return model{
		storage: store,
		player:  player,
		tracks:  t,
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
				}

				m.inputs[m.focused].Blur()
				m.nextInput()
				m.inputs[m.focused].Focus()
			}
		case "tab", "ctrl+n":
			switch m.view {
			case AddTrackView:
				m.inputs[m.focused].Blur()
				if m.view == AddTrackView {
					m.nextInput()
				}
				m.inputs[m.focused].Focus()
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
			m.tracks.Blur()
			m.view = AddTrackView
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
		s.WriteString(baseStyle.Render(m.tracks.View()) + helpStyle.Render(trackListHelp))
	case AddTrackView:
		s.WriteString(baseStyle.Render(m.tracks.View()) + "\n" + baseStyle.Render(
			fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
				inputHeaderStyle.Width(20).Render("Автор"),
				inputStyle.Render(m.inputs[artist].View()),
				inputHeaderStyle.Width(20).Render("Название"),
				inputStyle.Render(m.inputs[title].View()),
				inputHeaderStyle.Width(10).Render("Жанр"),
				inputStyle.Render(m.inputs[genre].View()),
				inputHeaderStyle.Width(10).Render("Тип файла"),
				inputStyle.Render(m.inputs[ft].View()),
				inputHeaderStyle.Width(30).Render("Путь к файлу"),
				inputStyle.Render(m.inputs[fp].View()),
			)) + helpStyle.Render(inputHelp),
		)
	case PlayerView:
	}

	return s.String()
}

func newTracksTable(store *storage.Storage) table.Model {
	tracks, _ := store.GetAllTracks()

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Исполнитель", Width: 12},
		{Title: "Название", Width: 12},
		{Title: "Тип файла", Width: 10},
		{Title: "Жанр", Width: 10},
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
		table.WithHeight(len(rows)+1),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(gloss.NormalBorder()).
		BorderForeground(gloss.Color("7")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(gloss.Color("7")).
		Background(gloss.Color("#306844")).
		Bold(true)
	t.SetStyles(s)

	return t
}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) quitInput() {
	m.tracks = newTracksTable(m.storage)
	m.focused = 0
	m.view = TrackListView
	m.tracks.Focus()
}

func (m *model) addTrack() {
	newTrack := music.New(
		m.inputs[title].Value(),
		m.inputs[artist].Value(),
		m.inputs[genre].Value(),
		m.inputs[ft].Value(),
		m.inputs[fp].Value(),
	)

	newTrack.Normalize()

	if err := newTrack.Validate(); err != nil {
		log.Fatal(err)
		m.quitInput()
	}

	m.storage.AddTrack(&newTrack)
	m.quitInput()
}
