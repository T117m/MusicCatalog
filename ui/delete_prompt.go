package ui

import (
	"strings"

	gloss "github.com/charmbracelet/lipgloss"
)

func (m *model) renderDeletePrompt() string {
	var sb strings.Builder

	sb.WriteString(gloss.PlaceHorizontal(40, gloss.Center, "Удаление трека"))

	if m.errMsg != nil {
		sb.WriteString("\n\nОшибка: ")
		sb.WriteString(m.errMsg.Error())

		return deleteStyle.Width(40).Render(sb.String())
	}

	sb.WriteString("\n\nВы уверены что хотите удалить этот трек?\n\n")
	sb.WriteString(m.tracks.SelectedRow()[title+1])
	sb.WriteString(" - ")
	sb.WriteString(m.tracks.SelectedRow()[artist+1])
	sb.WriteString(gloss.PlaceHorizontal(40, gloss.Center, "\n\ny - Да | n - Нет"))

	return deleteStyle.Width(40).Render(sb.String())
}
