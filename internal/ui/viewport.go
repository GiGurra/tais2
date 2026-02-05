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

	content := fb.String()

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(vpW).
		Height(vpH)

	return style.Render(content)
}
