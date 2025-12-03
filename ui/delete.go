package ui

import (
	"strings"
	//tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

// func (m *model) updateDeleteView(msg tea.Msg) (tea.Model, tea.Cmd) {
func (m *model) renderDeletePrompt() string {
	var sb strings.Builder

	sb.WriteString(gloss.PlaceHorizontal(40, gloss.Center, "Удаление трека"))
	sb.WriteString("\n\nВы уверены что хотите удалить этот трек?\n\n")
	sb.WriteString(m.tracks.SelectedRow()[1])
	sb.WriteString(" - ")
	sb.WriteString(m.tracks.SelectedRow()[2])
	sb.WriteString(gloss.PlaceHorizontal(40, gloss.Center, "\n\ny - Да | n - Нет"))

	return deleteStyle.Width(40).Render(sb.String())
}
