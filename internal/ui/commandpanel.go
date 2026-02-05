package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (g GameView) renderCommandPanel(width, height int) string {
	innerW := width - 2
	innerH := height - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	_, hasSelection := g.selectedIdx()

	var content string
	switch {
	case g.mode == modeMove:
		content = "[M] Move: click target\n[Esc] Cancel"
	case hasSelection:
		content = "[M] Move"
	default:
		content = "Select a unit"
	}

	content = lipgloss.Place(innerW, innerH, lipgloss.Center, lipgloss.Center, content)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(innerW).
		Height(innerH)

	return style.Render(content)
}
