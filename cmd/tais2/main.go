package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/GiGurra/tais2/internal/ui"
)

func main() {
	p := tea.NewProgram(
		ui.NewMainMenu(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
