package ui

import (
	"fmt"

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
	inputs[title].CharLimit = 20
	inputs[title].Width = 20
	inputs[title].Prompt = ""

	inputs[artist] = ti.New()
	inputs[artist].CharLimit = 20
	inputs[artist].Width = 20
	inputs[artist].Prompt = ""

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

	return inputs
}

func renderInputForm(inputs []ti.Model) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
		inputHeaderStyle.Width(20).Render("Название"),
		inputStyle.Render(inputs[title].View()),
		inputHeaderStyle.Width(20).Render("Автор"),
		inputStyle.Render(inputs[artist].View()),
		inputHeaderStyle.Width(10).Render("Жанр"),
		inputStyle.Render(inputs[genre].View()),
		inputHeaderStyle.Width(10).Render("Тип файла"),
		inputStyle.Render(inputs[ft].View()),
		inputHeaderStyle.Width(30).Render("Путь к файлу"),
		inputStyle.Render(inputs[fp].View()),
	)
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
	m.tracks = newTracksTable(m.storage)
	m.focused = 0
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

	m.focused = 0
	m.inputs[m.focused].Focus()
}
