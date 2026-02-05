package ui

import (
	"strings"

	"github.com/GiGurra/tais2/internal/game"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	minWidth  = 100
	minHeight = 40
)

type menuItem int

const (
	menuCampaign menuItem = iota
	menuSkirmish
	menuExit
	menuItemCount
)

func (m menuItem) label() string {
	switch m {
	case menuCampaign:
		return "CAMPAIGN"
	case menuSkirmish:
		return "SINGLE BATTLE"
	case menuExit:
		return "EXIT"
	default:
		return ""
	}
}

type MainMenu struct {
	width    int
	height   int
	selected menuItem
	tooSmall bool
}

func NewMainMenu() MainMenu {
	return MainMenu{}
}

func NewMainMenuWithSize(width, height int) MainMenu {
	return MainMenu{
		width:    width,
		height:   height,
		tooSmall: width < minWidth || height < minHeight,
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tooSmall = msg.Width < minWidth || msg.Height < minHeight

	case tea.KeyMsg:
		if m.tooSmall {
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.selected--
			if m.selected < 0 {
				m.selected = menuItemCount - 1
			}
		case "down", "j":
			m.selected++
			if m.selected >= menuItemCount {
				m.selected = 0
			}
		case "enter", " ":
			return m.activate()
		}

	case tea.MouseMsg:
		if m.tooSmall {
			return m, nil
		}
		// Mouse handling for menu items is done via click detection in view coordinates
		if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
			if item, ok := m.hitTestMenu(msg.Y); ok {
				m.selected = item
				return m.activate()
			}
		}
		if msg.Action == tea.MouseActionMotion || msg.Action == tea.MouseActionPress {
			if item, ok := m.hitTestMenu(msg.Y); ok {
				m.selected = item
			}
		}
	}

	return m, nil
}

func (m MainMenu) activate() (tea.Model, tea.Cmd) {
	switch m.selected {
	case menuSkirmish:
		scenario := game.NewScenario(128, 128)
		return NewGameView(&scenario, m.width, m.height), nil
	case menuExit:
		return m, tea.Quit
	default:
		return m, nil
	}
}

func (m MainMenu) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	if m.tooSmall {
		return m.renderTooSmall()
	}

	return m.renderMenu()
}

func (m MainMenu) renderTooSmall() string {
	msg := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")).
		Render("Please increase terminal size to at least 100x40")

	current := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render(strings.Repeat(" ", 4) + "current: " + itoa(m.width) + "x" + itoa(m.height))

	block := lipgloss.JoinVertical(lipgloss.Center, msg, current)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

// menuStartY returns the Y coordinate where the first menu button starts rendering.
func (m MainMenu) menuStartY() int {
	// Title art is 7 lines + 1 blank line gap = 8 lines before buttons.
	// Each button is 5 lines tall (1 padding + 1 content + 1 padding + 2 border) with 1 line gap between.
	titleLines := 8
	btnHeight := 5
	menuBlockHeight := int(menuItemCount)*btnHeight + (int(menuItemCount)-1)*1
	totalHeight := titleLines + menuBlockHeight
	return (m.height - totalHeight) / 2 + titleLines
}

func (m MainMenu) hitTestMenu(y int) (menuItem, bool) {
	startY := m.menuStartY()
	btnHeight := 5
	for i := 0; i < int(menuItemCount); i++ {
		// Each button occupies 5 rows, with 1 row gap between buttons
		btnTop := startY + i*(btnHeight+1)
		btnBottom := btnTop + btnHeight - 1
		if y >= btnTop && y <= btnBottom {
			return menuItem(i), true
		}
	}
	return 0, false
}

const titleArt = `
████████╗ █████╗ ██╗███████╗    ██████╗
╚══██╔══╝██╔══██╗██║██╔════╝    ╚════██╗
   ██║   ███████║██║███████╗     █████╔╝
   ██║   ██╔══██║██║╚════██║    ██╔═══╝
   ██║   ██║  ██║██║███████║    ███████╗
   ╚═╝   ╚═╝  ╚═╝╚═╝╚══════╝    ╚══════╝`

func (m MainMenu) renderMenu() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("2"))

	title := titleStyle.Render(titleArt)

	var buttons []string
	for i := 0; i < int(menuItemCount); i++ {
		item := menuItem(i)
		buttons = append(buttons, m.renderButton(item.label(), item == m.selected))
	}

	menu := lipgloss.JoinVertical(lipgloss.Center, buttons...)

	block := lipgloss.JoinVertical(lipgloss.Center,
		title,
		"",
		menu,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

func (m MainMenu) renderButton(label string, selected bool) string {
	width := 44

	base := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())

	if selected {
		return base.
			BorderForeground(lipgloss.Color("2")).
			Foreground(lipgloss.Color("2")).
			Bold(true).
			Render(label)
	}

	return base.
		BorderForeground(lipgloss.Color("8")).
		Foreground(lipgloss.Color("7")).
		Render(label)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}
