package ui

import (
	"strings"

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
	case menuExit:
		return m, tea.Quit
	default:
		// Other menu items are stubs for now
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
	// The menu block is vertically centered. Each button is 3 lines tall, with 1 line gap between.
	// Total menu height: menuItemCount*3 + (menuItemCount-1)*1 + title(1) + gap(2)
	titleLines := 4 // title + blank lines above/below
	menuBlockHeight := int(menuItemCount)*3 + (int(menuItemCount)-1)*1
	totalHeight := titleLines + menuBlockHeight
	return (m.height - totalHeight) / 2 + titleLines
}

func (m MainMenu) hitTestMenu(y int) (menuItem, bool) {
	startY := m.menuStartY()
	for i := 0; i < int(menuItemCount); i++ {
		// Each button occupies 3 rows, with 1 row gap between buttons
		btnTop := startY + i*4
		btnBottom := btnTop + 2
		if y >= btnTop && y <= btnBottom {
			return menuItem(i), true
		}
	}
	return 0, false
}

func (m MainMenu) renderMenu() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("2"))

	title := titleStyle.Render("T A I S  2")

	var buttons []string
	for i := 0; i < int(menuItemCount); i++ {
		item := menuItem(i)
		buttons = append(buttons, m.renderButton(item.label(), item == m.selected))
	}

	menu := lipgloss.JoinVertical(lipgloss.Center, buttons...)

	block := lipgloss.JoinVertical(lipgloss.Center,
		"",
		title,
		"",
		menu,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, block)
}

func (m MainMenu) renderButton(label string, selected bool) string {
	width := 30

	base := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Padding(0, 2).
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
