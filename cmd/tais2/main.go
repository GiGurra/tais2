package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/tais2/internal/ui"
	"github.com/spf13/cobra"
)

type Config struct {
	// Future CLI params go here
}

func main() {
	boa.CmdT[Config]{
		Use:   "tais2",
		Short: "Terminal-based RTS game",
		Long:  "tais2 - A terminal-based real-time strategy game with mouse support",
		ParamEnrich: boa.ParamEnricherCombine(
			boa.ParamEnricherName,
			boa.ParamEnricherShort,
			boa.ParamEnricherBool,
		),
		RunFunc: func(cfg *Config, cmd *cobra.Command, args []string) {
			p := tea.NewProgram(
				ui.NewMainMenu(),
				tea.WithAltScreen(),
				tea.WithMouseAllMotion(),
			)

			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		},
	}.Run()
}
