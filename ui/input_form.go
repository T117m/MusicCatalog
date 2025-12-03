package ui

import (
	"strings"

	"github.com/T117m/MusicCatalog/music"

	gloss "github.com/charmbracelet/lipgloss"
	ti "github.com/charmbracelet/bubbles/textinput"
)

const (
	title = iota
	artist
	genre
	ft
	fp
)

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

func (m *model) renderInputForm() string {
	var (
		formHeader string
		sb strings.Builder

		fieldHeaders = [5]string{"Название", "Исполнитель", "Жанр", "Тип файла", "Путь к файлу"}

		titleErr    = ""
		artistErr   = ""
		genreErr    = ""
		fileTypeErr = ""
		filePathErr = ""
	)

	switch m.view {
	case AddTrackView:
		formHeader = "Добавление трека\n"
	case EditTrackView:
		formHeader = "Редактирование трека\n"
	}

	if m.errMsg != nil {
		switch m.errMsg {
		case music.ErrEmptyTitle:
			titleErr = "! Название не может быть пустым!"
		case music.ErrEmptyArtist:
			artistErr = "! Поле автора не может быть пустым!"
		case music.ErrEmptyFileType:
			fileTypeErr = "! Тип файла не может быть пустым!"
		case music.ErrEmptyFilePath:
			filePathErr = "! Путь к файлу не может быть пустым!"
		case music.ErrUnsupportedFormat:
			fileTypeErr = "! Неподдерживаемый тип файла!"
			filePathErr = "! Возможно указан неправильный путь!"
		}
	}

	errs := [5]string{titleErr, artistErr, genreErr, fileTypeErr, filePathErr}

	sb.WriteString(gloss.PlaceHorizontal(30, gloss.Center, inputHeaderStyle.Render(formHeader)))

	for i, input := range m.inputs {
		writeInputField(&sb, fieldHeaders[i], errs[i], &input)
	}

	return sb.String()
}

func writeInputField(sb *strings.Builder, header, err string, input *ti.Model) {
	sb.WriteString("\n")
	sb.WriteString(inputHeaderStyle.Render(header))
	sb.WriteString(errorStyle.Render(err))
	sb.WriteString("\n")
	sb.WriteString(inputStyle.Render(input.View()))
}

func (m *model) nextInput() {
	m.inputs[m.focused].Blur()
	m.focused = (m.focused + 1) % len(m.inputs)
	m.inputs[m.focused].Focus()
}

func (m *model) prevInput() {
	m.focused--
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *model) quitInput() {
	m.resetInputs()
	m.inputs[m.focused].Blur()
	m.errMsg = nil
	m.view = TrackListView
	m.tracks.Focus()
}

func (m *model) setFocus(index int) {
	if index < 0 || index >= len(m.inputs) {
		return
	}

	m.inputs[m.focused].Blur()
	m.focused = index
	m.inputs[m.focused].Focus()
}

func (m *model) getInputs() (string, string, string, string, string) {
	return m.inputs[title].Value(), m.inputs[artist].Value(), m.inputs[genre].Value(),
		m.inputs[ft].Value(), m.inputs[fp].Value()
}

func (m *model) resetInputs() {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}

	m.setFocus(0)
}
