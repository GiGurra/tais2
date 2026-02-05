package ui

import (
	"github.com/GiGurra/tais2/internal/game"
	"github.com/charmbracelet/lipgloss"
)

func unitTypeName(ut game.UnitType) string {
	switch ut {
	case game.UnitPeasant:
		return "Peasant"
	case game.UnitFootman:
		return "Footman"
	case game.UnitTownHall:
		return "Town Hall"
	case game.UnitBarracks:
		return "Barracks"
	case game.UnitFarm:
		return "Farm"
	default:
		return "Unknown"
	}
}

func (g GameView) renderSelection(width, height int) string {
	innerW := width - 2
	innerH := height - 2
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	content := "No selection"

	if idx, ok := g.selectedIdx(); ok {
		w := &g.scenario.World
		e := &w.Entities[idx]

		lines := ""

		if e.Mask&game.MaskOwner != 0 {
			pid := w.Owner[idx].PlayerID
			color := lipgloss.Color(itoa(int(playerColor(pid))))
			tag := lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render("Player " + itoa(int(pid)))
			lines += tag + "  "
		}
		if e.Mask&game.MaskUnitType != 0 {
			lines += unitTypeName(w.UnitTypeComp[idx].Type)
		}
		if e.Mask&game.MaskHealth != 0 {
			hp := w.Health[idx]
			lines += "  HP: " + itoa(int(hp.Current)) + "/" + itoa(int(hp.Max))
		}
		if e.Mask&game.MaskPosition != 0 {
			pos := w.Position[idx]
			lines += "\nPos: " + itoa(int(pos.X/1000)) + "," + itoa(int(pos.Y/1000))
		}
		if e.Mask&game.MaskMovement != 0 {
			lines += "  Spd: " + itoa(int(w.Movement[idx].Speed))
		}
		if e.Mask&game.MaskMoveTarget != 0 {
			mt := w.MoveTarget[idx]
			lines += "\nMove: " + itoa(int(mt.X/1000)) + "," + itoa(int(mt.Y/1000))
		}

		content = lines
	}

	content = lipgloss.Place(innerW, innerH, lipgloss.Center, lipgloss.Center, content)

	style := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8")).
		Width(innerW).
		Height(innerH)

	return style.Render(content)
}
