package ui

import (
	"strings"

	"github.com/T117m/MusicCatalog/storage"

	"github.com/charmbracelet/bubbles/table"
	ti "github.com/charmbracelet/bubbles/textinput"
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

	input.Width = 40
	input.Prompt = ""

	return searchModel{
		input:   input,
		tag:     title,
		results: results,
	}
}

func (s *searchModel) View() string {
	var (
		sb         strings.Builder
		choosenTag = ""
	)

	switch s.tag {
	case title:
		choosenTag = "Название"
	case artist:
		choosenTag = "Автор"
	case genre:
		choosenTag = "Жанр"
	case ft:
		choosenTag = "Тип файла"
	}

	sb.WriteString(
		gloss.JoinHorizontal(
			gloss.Top,
			baseStyle.Width(40).Render(s.input.View()),
			searchFieldStyle.Width(9).AlignHorizontal(
				gloss.Center,
			).Render(choosenTag),
		),
	)

	return sb.String()
}
