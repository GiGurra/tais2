package main

import (
	"fmt"
	"os"

	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/tais2/internal/game"
	"github.com/GiGurra/tais2/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Config struct {
	// Future CLI params go here
}

type SingleBattleConfig struct {
	Snapshot bool `descr:"Print one frame and exit"`
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
			runInteractive(ui.NewMainMenu())
		},
		SubCmds: boa.SubCmds(
			boa.CmdT[SingleBattleConfig]{
				Use:   "single-battle",
				Short: "Jump straight into a single battle",
				ParamEnrich: boa.ParamEnricherCombine(
					boa.ParamEnricherName,
					boa.ParamEnricherShort,
					boa.ParamEnricherBool,
				),
				RunFunc: func(cfg *SingleBattleConfig, cmd *cobra.Command, args []string) {
					scenario := game.NewScenario(128, 128)
					game.SetupMatch(&scenario)
					if cfg.Snapshot {
						view := ui.NewGameView(&scenario, minWidth, minHeight)
						fmt.Print(view.View())
						return
					}
					runInteractive(ui.NewGameView(&scenario, 0, 0))
				},
			},
		),
	}.Run()
}

const (
	minWidth  = 100
	minHeight = 40
)

func runInteractive(model tea.Model) {
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
