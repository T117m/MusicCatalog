package ui

import (
	"github.com/charmbracelet/bubbles/table"
	gloss "github.com/charmbracelet/lipgloss"
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
	trackListHelp = "\nq: Выйти | a/i: Добавить трек | d: Удалить трек | r: Редактировать трек\n"
	inputHelp     = "\nq: Вернуться | Enter: Ввод\n"
)

func newStyledTable(columns []table.Column, rows []table.Row) table.Model {
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
