package ui

import (
	"github.com/GiGurra/tais2/internal/game"
	"github.com/charmbracelet/lipgloss"
)

// terrainGlyph maps a terrain type to a display rune, whether it's wide, and colors.
// Wide glyphs (like emoji) naturally occupy 2 terminal columns.
// Narrow glyphs are repeated to fill the 2-column tile.
func terrainGlyph(t game.TerrainType) (glyph rune, wide bool, fg, bg int32) {
	switch t {
	case game.Water:
		return '~', false, 34, 17 // blue fg, dark blue bg
	case game.Forest:
		return '🌲', true, 34, 22 // green fg, dark green bg
	case game.Dirt:
		return '.', false, 178, 94 // yellow fg, brown bg
	case game.Mountain:
		return '^', false, 255, 240 // white fg, gray bg
	case game.GoldMine:
		return '$', false, 220, 94 // yellow fg, brown bg
	default: // Grass
		return '.', false, 34, 22 // green fg, dark green bg
	}
}

func (g GameView) renderViewport() string {
	vpW, vpH := g.viewportSize()
	tilesW := vpW / 2

	fb := NewFrameBuffer(vpW, vpH)
	fb.Clear()

	// Blit terrain into framebuffer — each tile occupies 2 terminal columns
	for y := 0; y < vpH; y++ {
		for tx := 0; tx < tilesW; tx++ {
			tileX := int32(g.camX + tx)
			tileY := int32(g.camY + y)
			screenX := tx * 2
			if g.scenario.Terrain.InBounds(tileX, tileY) {
				tt := g.scenario.Terrain.At(tileX, tileY)
				glyph, wide, fg, bg := terrainGlyph(tt)
				if wide {
					fb.SetWide(screenX, y, glyph, fg, bg)
				} else {
					fb.Set(screenX, y, glyph, fg, bg)
					fb.Set(screenX+1, y, glyph, fg, bg)
				}
			}
		}
	}

	// Unit rendering pass — draw entities on top of terrain
	g.scenario.World.Each(game.MaskPosition|game.MaskRenderable, func(idx int32) bool {
		pos := g.scenario.World.Position[idx]
		tileX := int(pos.X / 1000)
		tileY := int(pos.Y / 1000)

		// Convert to screen coords relative to camera
		screenX := (tileX - g.camX) * 2
		screenY := tileY - g.camY

		// Skip off-screen entities
		if screenX < 0 || screenX+1 >= vpW || screenY < 0 || screenY >= vpH {
			return true
		}

		glyph := g.scenario.World.Renderable[idx].Glyph

		// Determine fg color: use player color for owned entities
		fg := g.scenario.World.Renderable[idx].Color
		if g.scenario.World.Entities[idx].Mask&game.MaskOwner != 0 {
			fg = playerColor(g.scenario.World.Owner[idx].PlayerID)
		}

		// Use terrain bg color for blending
		terrainTileX := int32(tileX)
		terrainTileY := int32(tileY)
		bg := int32(0)
		if g.scenario.Terrain.InBounds(terrainTileX, terrainTileY) {
			_, _, _, bg = terrainGlyph(g.scenario.Terrain.At(terrainTileX, terrainTileY))
		}

		fb.Set(screenX, screenY, glyph, fg, bg)
		fb.Set(screenX+1, screenY, glyph, fg, bg)

		return true
	})

	content := fb.String()

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(vpW).
		Height(vpH)

	return style.Render(content)
}

// playerColor returns the ANSI color for a player ID.
func playerColor(playerID int32) int32 {
	switch playerID {
	case 0:
		return 27 // blue
	case 1:
		return 196 // red
	default:
		return 7 // white
	}
}
