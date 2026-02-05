package ui

import (
	"strings"

	"github.com/GiGurra/tais2/internal/game"
	"github.com/charmbracelet/lipgloss"
)

func (g GameView) renderMinimap(width, height int) string {
	// Account for border (2 chars each side)
	innerW := width - 2
	innerH := height - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	var content string
	if g.scenario.Terrain.Width == 0 || g.scenario.Terrain.Height == 0 {
		content = lipgloss.Place(innerW, innerH, lipgloss.Center, lipgloss.Center, "No map")
	} else {
		// Scaled-down terrain representation
		var rows []string
		for y := 0; y < innerH; y++ {
			ty := int32(y) * g.scenario.Terrain.Height / int32(innerH)
			var row strings.Builder
			for x := 0; x < innerW; x++ {
				tx := int32(x) * g.scenario.Terrain.Width / int32(innerW)
				row.WriteByte(terrainChar(g.scenario.Terrain.At(tx, ty)))
			}
			rows = append(rows, row.String())
		}
		content = strings.Join(rows, "\n")
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(innerW).
		Height(innerH)

	return style.Render(content)
}

func terrainChar(t game.TerrainType) byte {
	switch t {
	case game.Water:
		return '~'
	case game.Forest:
		return 'T'
	case game.Dirt:
		return '.'
	case game.Mountain:
		return '^'
	default: // Grass
		return ' '
	}
}
