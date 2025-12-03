package ui

import (
	"github.com/charmbracelet/bubbles/table"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
	baseStyle = gloss.NewStyle().
			BorderStyle(gloss.ThickBorder()).
			BorderForeground(gloss.Color("7")).
			Foreground(gloss.Color("7"))
	deleteStyle = gloss.NewStyle().
			Foreground(gloss.Color("#FF746C")).
			BorderStyle(gloss.ThickBorder()).
			BorderForeground(gloss.Color("#FF746C")).
			Bold(true)
	errorStyle = gloss.NewStyle().
			Foreground(gloss.Color("#FF746C")).
			Bold(true)
	helpStyle = gloss.NewStyle().
			Foreground(gloss.Color("241"))
	inputHeaderStyle = gloss.NewStyle().
				Foreground(gloss.Color("7")).
				Bold(true)
	inputStyle = gloss.NewStyle().
			BorderStyle(gloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(gloss.Color("7"))
)

const (
	trackListHelp = "\nq: Выйти | Ctrl+a: Добавить трек | x: Удалить трек\nCtrl+r: Редактировать трек | Enter: Включить/Выключить трек\n"
	inputHelp     = "\nEsc: Вернуться | Enter: Ввод | Ctrl+s: Сохранить\n"
	deleteHelp    = "\nEsc/q: Вернутся\n"
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
		BorderBottom(true)
	s.Selected = s.Selected.
		Background(gloss.Color("#306844")).
		Foreground(gloss.Color("7"))
	t.SetStyles(s)

	return t
}
