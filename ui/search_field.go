package ui

import (
	"strings"

	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	ti "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type searchModel struct {
	input   ti.Model
	tag     field
	results table.Model
}

func newSearchModel(store *storage.Storage) searchModel {
	var (
		input   = ti.New()
		results = newTrackList(store)
	)

	input.Width = 45
	input.Prompt = ""

	return searchModel{
		input:   input,
		tag:     title,
		results: results,
	}
}

func (s searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

func (s *searchModel) View() string {
	var (
		sb     strings.Builder
		titles = [4]string{"Название", "Автор", "Жанр", "Тип файла"}
	)

	if s.tag > 3 || s.tag < 0 {
		s.tag = 0
	}

	sb.WriteString(
		gloss.JoinHorizontal(
			gloss.Top,
			baseStyle.Width(46).Render(s.input.View()),
			searchFieldStyle.Width(10).AlignHorizontal(
				gloss.Center,
			).Render(titles[s.tag]),
		),
	)
	sb.WriteString("\n" + baseStyle.Render(s.results.View()))

	return sb.String()
}
