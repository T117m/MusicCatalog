package ui

import (
	"strings"

	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type searchModel struct {
	input   ti.Model
	tag     field
}

func newSearchModel() searchModel {
	input := ti.New()

	input.Width = 45
	input.Prompt = ""

	return searchModel{
		input: input,
		tag:   title,
	}
}

func (s searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

var titles = [4]string{"Название", "Автор", "Жанр", "Тип файла"}

func (s *searchModel) View(st gloss.Style) string {
	var (
		sb strings.Builder
	)

	sb.WriteString(
		gloss.JoinHorizontal(
			gloss.Top,
			st.Width(46).Render(s.input.View()),
			st.Width(10).AlignHorizontal(
				gloss.Center,
			).Render(titles[s.tag]),
		),
	)

	return sb.String()
}

func (s *searchModel) nextTag() {
	if s.tag == ft {
		s.tag = title
	} else {
		s.tag++
	}
}

func (s *searchModel) prevTag() {
	if s.tag == title {
		s.tag = ft
	} else {
		s.tag--
	}
}
