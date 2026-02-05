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

		// Build a 2D grid of characters and colors for the minimap.
		// Start with terrain, then stamp entities on top.
		type cell struct {
			ch byte
			fg int32 // 0 = default (no color override)
		}
		grid := make([]cell, innerW*innerH)
		for y := 0; y < innerH; y++ {
			ty := int32(y) * g.scenario.Terrain.Height / int32(innerH)
			for x := 0; x < innerW; x++ {
				tx := int32(x) * g.scenario.Terrain.Width / int32(innerW)
				grid[y*innerW+x] = cell{ch: terrainChar(g.scenario.Terrain.At(tx, ty))}
			}
		}

		// Draw entities on the minimap
		g.scenario.World.Each(game.MaskPosition|game.MaskOwner, func(idx int32) bool {
			pos := g.scenario.World.Position[idx]
			tileX := int(pos.X / 1000)
			tileY := int(pos.Y / 1000)

			mx := tileX * innerW / mapW
			my := tileY * innerH / mapH
			if mx >= 0 && mx < innerW && my >= 0 && my < innerH {
				pid := g.scenario.World.Owner[idx].PlayerID
				grid[my*innerW+mx] = cell{ch: '*', fg: playerColor(pid)}
			}
			return true
		})

		var rows []string
		for y := 0; y < innerH; y++ {
			var row strings.Builder
			x := 0
			for x < innerW {
				c := grid[y*innerW+x]
				inCam := x >= camLeft && x < camRight && y >= camTop && y < camBot

				// Group consecutive cells with the same styling
				var span strings.Builder
				hasFG := c.fg != 0
				fg := c.fg
				span.WriteByte(c.ch)
				x++
				for x < innerW {
					c2 := grid[y*innerW+x]
					inCam2 := x >= camLeft && x < camRight && y >= camTop && y < camBot
					hasFG2 := c2.fg != 0
					if inCam2 != inCam || hasFG2 != hasFG || (hasFG2 && c2.fg != fg) {
						break
					}
					span.WriteByte(c2.ch)
					x++
				}

				s := span.String()
				style := lipgloss.NewStyle()
				if hasFG {
					style = style.Foreground(lipgloss.Color(itoa(int(fg))))
				}
				if inCam {
					style = style.Background(lipgloss.Color("236"))
				}
				if hasFG || inCam {
					row.WriteString(style.Render(s))
				} else {
					row.WriteString(s)
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

// minimapScreenRect returns the screen bounds of the minimap inner area.
// Used for click detection.
func (g GameView) minimapScreenRect() (x0, y0, x1, y1 int) {
	_, vpH := g.viewportSize()
	minimapW, _, _ := g.panelWidths()
	innerW := minimapW - 2

	// Layout: HUD (1 row) + viewport border top (1) + vpH + viewport border bottom (1)
	// = 1 + 1 + vpH + 1 = vpH + 3
	topOfBottom := vpH + 3
	// Minimap border adds 1 row top, 1 col left
	x0 = 1            // left border of minimap
	y0 = topOfBottom + 1 // top border of minimap
	x1 = x0 + innerW
	y1 = y0 + (bottomPanelHeight - 2)
	return
}

// screenToMinimapTile converts a screen click inside the minimap to a map tile coordinate.
func (g GameView) screenToMinimapTile(mx, my int) (tileX, tileY int, ok bool) {
	x0, y0, x1, y1 := g.minimapScreenRect()
	if mx < x0 || mx >= x1 || my < y0 || my >= y1 {
		return 0, 0, false
	}

	innerW := x1 - x0
	innerH := y1 - y0
	relX := mx - x0
	relY := my - y0

	mapW := int(g.scenario.Terrain.Width)
	mapH := int(g.scenario.Terrain.Height)

	tileX = relX * mapW / innerW
	tileY = relY * mapH / innerH
	return tileX, tileY, true
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
	case game.GoldMine:
		return '$'
	default: // Grass
		return ' '
	}
}
