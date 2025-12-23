package ui

import (
	"errors"
	"strings"

	"github.com/T117m/MusicCatalog/music"
	"github.com/T117m/MusicCatalog/storage"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type field int

const (
	title field = iota
	artist
	genre
	ft
	fp
)

type inputFormModel struct {
	inputs  []ti.Model
	focused int
}

func newInputModel() inputFormModel {
	return inputFormModel{
		inputs:  newInputs(),
		focused: 0,
	}
}

func newInputs() []ti.Model {
	inputs := make([]ti.Model, 5)

	inputs[title] = ti.New()
	inputs[title].Width = 27
	inputs[title].Prompt = ""

	inputs[artist] = ti.New()
	inputs[artist].Width = 27
	inputs[artist].Prompt = ""

	inputs[genre] = ti.New()
	inputs[genre].Width = 27
	inputs[genre].Prompt = ""

	inputs[ft] = ti.New()
	inputs[ft].CharLimit = 4
	inputs[ft].Width = 4
	inputs[ft].Prompt = ""

	inputs[fp] = ti.New()
	inputs[fp].Width = 27
	inputs[fp].Prompt = ""

	return inputs
}

func (ifm inputFormModel) Update(msg tea.Msg) (inputFormModel, tea.Cmd) {
	var cmd tea.Cmd
	ifm.inputs[ifm.focused], cmd = ifm.inputs[ifm.focused].Update(msg)
	return ifm, cmd
}

func (ifm *inputFormModel) Focus() tea.Cmd {
	return ifm.inputs[ifm.focused].Focus()
}

func (ifm *inputFormModel) Blur() {
	ifm.inputs[ifm.focused].Blur()
}

func (m *model) writeInputForm(sb *strings.Builder, k keymap) {
	sb.WriteString(
		gloss.JoinHorizontal(
			gloss.Top,
			gloss.JoinVertical(
				gloss.Center,
				m.search.View(unfocusedStyle),
				unfocusedStyle.Render(m.tracks.View()),
			),
			baseStyle.Width(30).Render(m.renderInputForm()),
		))
	sb.WriteString(gloss.PlaceHorizontal(
		gloss.Width(baseStyle.Render(m.tracks.View()))+
			gloss.Width("\n" + m.help.View(k)),
		gloss.Right,
		helpStyle.Render("\n" + m.help.View(k)),
	))
}

func (m *model) renderInputForm() string {
	var (
		formHeader string
		sb         strings.Builder

		fieldHeaders = [5]string{"Название", "Исполнитель", "Жанр", "Тип файла", "Путь к файлу"}

		titleErr    = ""
		artistErr   = ""
		genreErr    = ""
		fileTypeErr = ""
		filePathErr = ""
		defaultErr  = ""
	)

	switch m.view {
	case AddTrackView:
		formHeader = "Добавление трека\n"
	case EditTrackView:
		formHeader = "Редактирование трека\n"
	}

	if m.errMsg != nil {
		switch {
		case errors.Is(m.errMsg, music.ErrEmptyTitle):
			titleErr = "! Название не может быть пустым!"
		case errors.Is(m.errMsg, music.ErrEmptyArtist):
			artistErr = "! Поле автора не может быть пустым!"
		case errors.Is(m.errMsg, music.ErrEmptyFileType):
			fileTypeErr = "! Тип файла не может быть пустым!"
		case errors.Is(m.errMsg, music.ErrEmptyFilePath):
			filePathErr = "! Путь к файлу не может быть пустым!"
		case errors.Is(m.errMsg, music.ErrUnsupportedFormat):
			fileTypeErr = "! Неподдерживаемый тип файла!"
			filePathErr = "! Возможно указан неправильный путь!"
		case errors.Is(m.errMsg, storage.ErrNotUnique):
			filePathErr = "! Этот файл уже храниться в каталоге!"
		case errors.Is(m.errMsg, storage.ErrFileNotExists):
			filePathErr = "! Файл не существует!"
		case errors.Is(m.errMsg, storage.ErrNoFileAccess):
			filePathErr = "! Не удаётся открыть файл!"
		case errors.Is(m.errMsg, storage.ErrIsDirectory):
			filePathErr = "! Не является файлом!"
		default:
			defaultErr = m.errMsg.Error()
		}
	}

	errs := [5]string{titleErr, artistErr, genreErr, fileTypeErr, filePathErr}

	sb.WriteString(gloss.PlaceHorizontal(30, gloss.Center, inputHeaderStyle.Render(formHeader)))

	for i, input := range m.input.inputs {
		writeInputField(&sb, fieldHeaders[i], errs[i], &input)
	}

	sb.WriteString("\n" + errorStyle.Render(defaultErr))

	return sb.String()
}

func writeInputField(sb *strings.Builder, header, errMsg string, input *ti.Model) {
	sb.WriteString("\n")
	sb.WriteString(inputHeaderStyle.Render(header))
	sb.WriteString(errorStyle.Render(errMsg))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(input.View()))
}

func (ifm *inputFormModel) nextInput() {
	ifm.Blur()
	ifm.focused = (ifm.focused + 1) % len(ifm.inputs)
	ifm.Focus()
}

func (ifm *inputFormModel) prevInput() {
	ifm.focused--
	if ifm.focused < 0 {
		ifm.focused = len(ifm.inputs) - 1
	}
}

func (m *model) quitInput() {
	m.input.resetInputs()
	m.input.Blur()
	m.errMsg = nil
	m.view = TrackListView
	m.tracks.Focus()
}

func (ifm *inputFormModel) setFocus(f field) {
	index := int(f)

	if index < 0 || index >= len(ifm.inputs) {
		return
	}

	ifm.Blur()
	ifm.focused = index
	ifm.Focus()
}

func (ifm *inputFormModel) getInputs() (string, string, string, string, string) {
	return ifm.inputs[title].Value(), ifm.inputs[artist].Value(), ifm.inputs[genre].Value(),
		ifm.inputs[ft].Value(), ifm.inputs[fp].Value()
}

func (ifm *inputFormModel) resetInputs() {
	for i := range ifm.inputs {
		ifm.inputs[i].Reset()
	}

	ifm.setFocus(0)
}
