package ui

import (
	"strings"
	"time"

	"github.com/GiGurra/tais2/internal/game"
	"github.com/GiGurra/tais2/internal/game/script"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const bottomPanelHeight = 10

type tickMsg struct{}

// inputMode tracks what happens on the next left-click.
type inputMode int

const (
	modeSelect inputMode = iota // default: left-click selects
	modeMove                    // left-click issues move order
)

type GameView struct {
	scenario   *game.Scenario
	width      int
	height     int
	tooSmall   bool
	camX       int // top-left tile X of viewport
	camY       int // top-left tile Y of viewport
	selection  game.EntityID    // selected entity (-1 = none)
	scriptExec *script.Executor // optional script executor
	speed      int32            // tick rate multiplier (1 = normal)
	mode       inputMode        // current input mode
	debug      bool             // show debug info in HUD
	lastMouseX int              // last mouse screen X (for debug)
	lastMouseY int              // last mouse screen Y (for debug)
}

func NewGameView(scenario *game.Scenario, width, height int) GameView {
	return GameView{
		scenario:  scenario,
		width:     width,
		height:    height,
		tooSmall:  width < minWidth || height < minHeight,
		selection: -1,
		speed:     1,
	}
}

func NewGameViewWithScript(scenario *game.Scenario, width, height int, exec *script.Executor, speed int32) GameView {
	return GameView{
		scenario:   scenario,
		width:      width,
		height:     height,
		tooSmall:   width < minWidth || height < minHeight,
		selection:  -1,
		scriptExec: exec,
		speed:      speed,
	}
}

// WithDebug returns a copy of the GameView with debug mode enabled.
func (g GameView) WithDebug(debug bool) GameView {
	g.debug = debug
	return g
}

func (g GameView) Init() tea.Cmd {
	return g.tickCmd()
}

func (g GameView) tickCmd() tea.Cmd {
	rate := g.scenario.TickRate
	if rate <= 0 {
		rate = 10
	}
	if g.speed > 1 {
		rate *= g.speed
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
		if g.scriptExec != nil {
			halt, _ := g.scriptExec.ExecTick(g.scenario)
			if halt {
				return g, nil // stop ticking
			}
		}
		g.scenario.Step()
		return g, g.tickCmd()

	case tea.MouseMsg:
		g.lastMouseX = msg.X
		g.lastMouseY = msg.Y

		// Minimap drag: press or motion with left button held
		if msg.Button == tea.MouseButtonLeft &&
			(msg.Action == tea.MouseActionPress || msg.Action == tea.MouseActionMotion) {
			if tileX, tileY, ok := g.screenToMinimapTile(msg.X, msg.Y); ok {
				g.centerCameraOn(tileX, tileY)
				return g, nil
			}
		}

		if msg.Action == tea.MouseActionRelease {
			switch msg.Button {
			case tea.MouseButtonLeft:
				if tileX, tileY, ok := g.screenToMinimapTile(msg.X, msg.Y); ok {
					g.centerCameraOn(tileX, tileY)
				} else if g.mode == modeMove {
					g.issueMoveClick(msg.X, msg.Y)
					g.mode = modeSelect
				} else {
					if tileX, tileY, ok := g.screenToTile(msg.X, msg.Y); ok {
						g.selection = g.findEntityAtTile(tileX, tileY)
					}
				}
			case tea.MouseButtonRight:
				// Right-click move (works on Alacritty, Kitty, WezTerm; not iTerm2)
				g.issueMoveClick(msg.X, msg.Y)
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return g, tea.Quit
		case "esc":
			if g.mode != modeSelect {
				g.mode = modeSelect
				return g, nil
			}
			return NewMainMenuWithSize(g.width, g.height), nil
		case "m":
			if idx, ok := g.selectedIdx(); ok {
				if g.scenario.World.Entities[idx].Mask&game.MaskMovement != 0 {
					g.mode = modeMove
					return g, nil
				}
			}
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

// centerCameraOn moves the camera so the given tile is at the center of the viewport.
func (g *GameView) centerCameraOn(tileX, tileY int) {
	vpW, vpH := g.viewportSize()
	tilesW := vpW / 2
	g.camX = tileX - tilesW/2
	g.camY = tileY - vpH/2
	g.clampCamera()
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

// issueMoveClick issues a move order for the selected unit to the clicked tile.
func (g *GameView) issueMoveClick(mx, my int) {
	selIdx, ok := g.selectedIdx()
	if !ok {
		return
	}
	tileX, tileY, ok := g.screenToTile(mx, my)
	if !ok {
		return
	}
	targetX := tileX*1000 + 500
	targetY := tileY*1000 + 500
	g.scenario.World.MoveTarget[selIdx] = game.MoveTarget{X: targetX, Y: targetY}
	g.scenario.World.Entities[selIdx].Mask |= game.MaskMoveTarget
}

// selectedIdx returns the slot index of the selected entity if it is alive.
func (g GameView) selectedIdx() (int32, bool) {
	if g.selection < 0 {
		return 0, false
	}
	if !g.scenario.World.IsAlive(g.selection) {
		return 0, false
	}
	return g.selection.Index(), true
}

// screenToTile converts screen coordinates to tile coordinates, accounting for
// border (1 col each side) and HUD (1 row). Returns ok=false if outside viewport.
func (g GameView) screenToTile(mx, my int) (tileX, tileY int32, ok bool) {
	relX := mx - 1  // 1 for left border
	relY := my - 2  // 1 for HUD + 1 for top border
	vpW, vpH := g.viewportSize()
	if relX < 0 || relX >= vpW || relY < 0 || relY >= vpH {
		return 0, 0, false
	}
	tileX = int32(g.camX + relX/2)
	tileY = int32(g.camY + relY)
	return tileX, tileY, true
}

// findEntityAtTile returns the EntityID of a player-0 entity at the given tile, or -1.
// Only units owned by the local player (0) can be selected.
func (g GameView) findEntityAtTile(tileX, tileY int32) game.EntityID {
	result := game.EntityID(-1)
	g.scenario.World.Each(game.MaskPosition|game.MaskOwner, func(idx int32) bool {
		if g.scenario.World.Owner[idx].PlayerID != 0 {
			return true
		}
		pos := g.scenario.World.Position[idx]
		etx := pos.X / 1000
		ety := pos.Y / 1000
		if etx == tileX && ety == tileY {
			e := &g.scenario.World.Entities[idx]
			result = game.MakeEntityID(idx, e.Gen)
			return false // stop iteration
		}
		return true
	})
	return result
}

// panelWidths returns the widths for minimap (20%), selection (50%), commands (30%).
func (g GameView) panelWidths() (minimap, selection, commands int) {
	w := g.width
	minimap = w * 20 / 100
	commands = w * 30 / 100
	selection = w - minimap - commands // remainder to avoid off-by-one
	return
}
