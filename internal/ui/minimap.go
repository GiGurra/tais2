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
		mapW := int(g.scenario.Terrain.Width)
		mapH := int(g.scenario.Terrain.Height)
		vpW, vpH := g.viewportSize()
		tilesW := vpW / 2 // each tile is 2 terminal columns wide

		// Camera rect in minimap coordinates
		camLeft := g.camX * innerW / mapW
		camTop := g.camY * innerH / mapH
		camRight := (g.camX + tilesW) * innerW / mapW
		camBot := (g.camY + vpH) * innerH / mapH

		camBG := lipgloss.NewStyle().Background(lipgloss.Color("236"))

		var rows []string
		for y := 0; y < innerH; y++ {
			ty := int32(y) * g.scenario.Terrain.Height / int32(innerH)
			var row strings.Builder
			x := 0
			for x < innerW {
				tx := int32(x) * g.scenario.Terrain.Width / int32(innerW)
				ch := string(terrainChar(g.scenario.Terrain.At(tx, ty)))
				inCam := x >= camLeft && x < camRight && y >= camTop && y < camBot

				// Group consecutive cells with the same in/out-of-camera state
				var span strings.Builder
				span.WriteString(ch)
				x++
				for x < innerW {
					tx2 := int32(x) * g.scenario.Terrain.Width / int32(innerW)
					ch2 := string(terrainChar(g.scenario.Terrain.At(tx2, ty)))
					inCam2 := x >= camLeft && x < camRight && y >= camTop && y < camBot
					if inCam2 != inCam {
						break
					}
					span.WriteString(ch2)
					x++
				}

				if inCam {
					row.WriteString(camBG.Render(span.String()))
				} else {
					row.WriteString(span.String())
				}
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
