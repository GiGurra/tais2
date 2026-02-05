package ui

import (
	"strings"
	"time"

	"github.com/GiGurra/tais2/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const bottomPanelHeight = 10

type tickMsg struct{}

type GameView struct {
	scenario *game.Scenario
	width    int
	height   int
	tooSmall bool
	camX     int // top-left tile X of viewport
	camY     int // top-left tile Y of viewport
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
	return g.tickCmd()
}

func (g GameView) tickCmd() tea.Cmd {
	rate := g.scenario.TickRate
	if rate <= 0 {
		rate = 10
	}
	return tea.Tick(time.Second/time.Duration(rate), func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (g GameView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		g.width = msg.Width
		g.height = msg.Height
		g.tooSmall = msg.Width < minWidth || msg.Height < minHeight
		g.clampCamera()

	case tickMsg:
		g.scenario.Step()
		return g, g.tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return g, tea.Quit
		case "esc":
			return NewMainMenuWithSize(g.width, g.height), nil
		case "up", "w":
			g.camY--
			g.clampCamera()
		case "down", "s":
			g.camY++
			g.clampCamera()
		case "left", "a":
			g.camX--
			g.clampCamera()
		case "right", "d":
			g.camX++
			g.clampCamera()
		case "ctrl+up", "alt+up", "shift+up":
			_, vpH := g.viewportSize()
			g.camY -= vpH
			g.clampCamera()
		case "ctrl+down", "alt+down", "shift+down":
			_, vpH := g.viewportSize()
			g.camY += vpH
			g.clampCamera()
		case "ctrl+left", "alt+left", "shift+left":
			vpW, _ := g.viewportSize()
			g.camX -= vpW / 2
			g.clampCamera()
		case "ctrl+right", "alt+right", "shift+right":
			vpW, _ := g.viewportSize()
			g.camX += vpW / 2
			g.clampCamera()
		}
	}

	return g, nil
}

func (g *GameView) clampCamera() {
	mapW := int(g.scenario.Terrain.Width)
	mapH := int(g.scenario.Terrain.Height)
	vpW, vpH := g.viewportSize()
	tilesW := vpW / 2 // each tile is 2 terminal columns wide

	maxX := mapW - tilesW
	maxY := mapH - vpH
	if maxX < 0 {
		maxX = 0
	}
	if maxY < 0 {
		maxY = 0
	}

	if g.camX < 0 {
		g.camX = 0
	}
	if g.camX > maxX {
		g.camX = maxX
	}
	if g.camY < 0 {
		g.camY = 0
	}
	if g.camY > maxY {
		g.camY = maxY
	}
}

// viewportSize returns the inner dimensions of the viewport area (excluding border).
func (g GameView) viewportSize() (w, h int) {
	// 1 row for HUD, bottomPanelHeight for bottom panel, 2 for border
	h = g.height - 1 - bottomPanelHeight - 2
	if h < 1 {
		h = 1
	}
	w = g.width - 2 // 2 for border
	if w < 1 {
		w = 1
	}
	return
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
