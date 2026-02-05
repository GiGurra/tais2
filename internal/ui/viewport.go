package ui

import (
	"github.com/charmbracelet/lipgloss"
)

func (g GameView) renderViewport() string {
	// 1 row for HUD, bottomPanelHeight for bottom panel, 2 for border
	vpHeight := g.height - 1 - bottomPanelHeight - 2
	if vpHeight < 1 {
		vpHeight = 1
	}
	vpWidth := g.width - 2 // 2 for border
	if vpWidth < 1 {
		vpWidth = 1
	}

	placeholder := "Empty battlefield"

	content := lipgloss.Place(vpWidth, vpHeight, lipgloss.Center, lipgloss.Center, placeholder)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(vpWidth).
		Height(vpHeight)

	return style.Render(content)
}
