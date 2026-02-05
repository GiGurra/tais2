package ui

import (
	"strings"

	"github.com/GiGurra/tais2/internal/game"
	"github.com/charmbracelet/lipgloss"
)

func (g GameView) renderHUD() string {
	gold := int32(0)
	lumber := int32(0)
	g.scenario.World.Each(game.MaskResourceStore, func(idx int32) bool {
		rs := &g.scenario.World.ResourceStore[idx]
		gold += rs.Gold
		lumber += rs.Lumber
		return true
	})

	left := "  Gold: " + itoa(int(gold)) +
		"    Lumber: " + itoa(int(lumber)) +
		"    Supply: 0/0" +
		"    Tick: " + itoa(int(g.scenario.Tick))

	right := "[ESC] Menu  "

	gap := g.width - len(left) - len(right)
	if gap < 1 {
		gap = 1
	}

	row := left + strings.Repeat(" ", gap) + right

	style := lipgloss.NewStyle().
		Width(g.width).
		Foreground(lipgloss.Color("7")).
		Background(lipgloss.Color("0"))

	return style.Render(row)
}
