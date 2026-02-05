package ui

import (
	"strings"

	"github.com/GiGurra/tais2/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const bottomPanelHeight = 10

type GameView struct {
	scenario *game.Scenario
	width    int
	height   int
	tooSmall bool
}

func NewGameView(scenario *game.Scenario, width, height int) GameView {
	return GameView{
		scenario: scenario,
		width:    width,
		height:   height,
		tooSmall: width < minWidth || height < minHeight,
	}
}

func (g GameView) Init() tea.Cmd {
	return nil
}

func (g GameView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
		g.tooSmall = msg.Width < minWidth || msg.Height < minHeight

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "esc":
			return NewMainMenuWithSize(g.width, g.height), nil
		}
	}

	return g, nil
}

func (g GameView) View() string {
	if g.width == 0 || g.height == 0 {
		return ""
	}

	if g.tooSmall {
		msg := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			Render("Please increase terminal size to at least 100x40")

		current := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render(strings.Repeat(" ", 4) + "current: " + itoa(g.width) + "x" + itoa(g.height))

		block := lipgloss.JoinVertical(lipgloss.Center, msg, current)
		return lipgloss.Place(g.width, g.height, lipgloss.Center, lipgloss.Center, block)
	}

	hud := g.renderHUD()
	viewport := g.renderViewport()
	bottom := g.renderBottomPanel()

	return lipgloss.JoinVertical(lipgloss.Left, hud, viewport, bottom)
}

func (g GameView) renderBottomPanel() string {
	minimapW, selectionW, commandsW := g.panelWidths()

	minimap := g.renderMinimap(minimapW, bottomPanelHeight)
	selection := g.renderSelection(selectionW, bottomPanelHeight)
	commands := g.renderCommandPanel(commandsW, bottomPanelHeight)

	return lipgloss.JoinHorizontal(lipgloss.Top, minimap, selection, commands)
}

// panelWidths returns the widths for minimap (20%), selection (50%), commands (30%).
func (g GameView) panelWidths() (minimap, selection, commands int) {
	w := g.width
	minimap = w * 20 / 100
	commands = w * 30 / 100
	selection = w - minimap - commands // remainder to avoid off-by-one
	return
}
