package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Cell represents a single character cell in the framebuffer.
// A Glyph of 0 marks a continuation cell (occupied by a preceding wide glyph).
type Cell struct {
	Glyph rune
	FG    int32 // ANSI color (0 = default)
	BG    int32 // ANSI color (0 = default)
}

// FrameBuffer is a rune+color grid sized to the viewport, rewritten each frame.
type FrameBuffer struct {
	Width, Height int
	Cells         []Cell // row-major: index = y*Width + x
}

func NewFrameBuffer(w, h int) FrameBuffer {
	return FrameBuffer{
		Width:  w,
		Height: h,
		Cells:  make([]Cell, w*h),
	}
}

func (fb *FrameBuffer) Resize(w, h int) {
	if w == fb.Width && h == fb.Height {
		return
	}
	fb.Width = w
	fb.Height = h
	fb.Cells = make([]Cell, w*h)
}

func (fb *FrameBuffer) Clear() {
	for i := range fb.Cells {
		fb.Cells[i] = Cell{Glyph: ' '}
	}
}

func (fb *FrameBuffer) Set(x, y int, glyph rune, fg, bg int32) {
	if x < 0 || x >= fb.Width || y < 0 || y >= fb.Height {
		return
	}
	fb.Cells[y*fb.Width+x] = Cell{Glyph: glyph, FG: fg, BG: bg}
}

// SetWide places a double-width glyph at (x,y) and marks (x+1,y) as a continuation cell.
func (fb *FrameBuffer) SetWide(x, y int, glyph rune, fg, bg int32) {
	fb.Set(x, y, glyph, fg, bg)
	if x+1 < fb.Width {
		fb.Cells[y*fb.Width+x+1] = Cell{Glyph: 0, FG: fg, BG: bg}
	}
}

func (fb *FrameBuffer) Get(x, y int) Cell {
	if x < 0 || x >= fb.Width || y < 0 || y >= fb.Height {
		return Cell{Glyph: ' '}
	}
	return fb.Cells[y*fb.Width+x]
}

// String renders the framebuffer to an ANSI string using lipgloss for styling.
// Consecutive cells with the same FG/BG are grouped into styled spans.
// Continuation cells (Glyph == 0) are skipped — their column is occupied
// by the preceding wide glyph.
func (fb *FrameBuffer) String() string {
	var rows []string
	for y := 0; y < fb.Height; y++ {
		var row strings.Builder
		x := 0
		for x < fb.Width {
			cell := fb.Cells[y*fb.Width+x]
			if cell.Glyph == 0 {
				x++
				continue // continuation cell
			}

			fg := cell.FG
			bg := cell.BG

			// Collect consecutive cells with the same colors
			var span strings.Builder
			span.WriteRune(cell.Glyph)
			x++

			for x < fb.Width {
				next := fb.Cells[y*fb.Width+x]
				if next.Glyph == 0 {
					x++ // skip continuation, same colors
					continue
				}
				if next.FG != fg || next.BG != bg {
					break
				}
				span.WriteRune(next.Glyph)
				x++
			}

			// Style the span
			style := lipgloss.NewStyle()
			if fg != 0 {
				style = style.Foreground(lipgloss.Color(fmt.Sprintf("%d", fg)))
			}
			if bg != 0 {
				style = style.Background(lipgloss.Color(fmt.Sprintf("%d", bg)))
			}
			row.WriteString(style.Render(span.String()))
		}
		rows = append(rows, row.String())
	}
	return strings.Join(rows, "\n")
}
