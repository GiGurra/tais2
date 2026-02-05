package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (g GameView) renderSelection(width, height int) string {
	innerW := width - 2
	innerH := height - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	content := lipgloss.Place(innerW, innerH, lipgloss.Center, lipgloss.Center, "No selection")

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(innerW).
		Height(innerH)

	return style.Render(content)
}
