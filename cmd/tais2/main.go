package main

import (
	"fmt"
	"os"

	"github.com/GiGurra/boa/pkg/boa"
	"github.com/GiGurra/tais2/internal/game"
	"github.com/GiGurra/tais2/internal/game/script"
	"github.com/GiGurra/tais2/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type Config struct {
	// Future CLI params go here
}

type SingleBattleConfig struct {
	Snapshot bool   `descr:"Print one frame and exit"`
	Script   string `descr:"Path to script file for playback"`
	Speed    int32  `descr:"Tick rate multiplier" default:"1"`
}

type SimulateConfig struct {
	Script   string `descr:"Path to script file" required:"true"`
	MaxTicks int32  `descr:"Maximum ticks to simulate (0=unlimited)" default:"10000"`
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
					if cfg.Script != "" {
						sc, exec := loadScript(cfg.Script)
						_ = sc
						speed := cfg.Speed
						if speed < 1 {
							speed = 1
						}
						runInteractive(ui.NewGameViewWithScript(&scenario, 0, 0, exec, speed))
						return
					}
					runInteractive(ui.NewGameView(&scenario, 0, 0))
				},
			},
			boa.CmdT[SimulateConfig]{
				Use:   "simulate",
				Short: "Run a headless simulation from a script file",
				ParamEnrich: boa.ParamEnricherCombine(
					boa.ParamEnricherName,
					boa.ParamEnricherShort,
					boa.ParamEnricherBool,
				),
				RunFunc: func(cfg *SimulateConfig, cmd *cobra.Command, args []string) {
					runSimulate(cfg)
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

func loadScript(path string) (*script.Script, *script.Executor) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening script: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	sc, err := script.Parse(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing script: %v\n", err)
		os.Exit(1)
	}

	exec := script.NewExecutor(sc, os.Stderr)
	return sc, exec
}

func runSimulate(cfg *SimulateConfig) {
	_, exec := loadScript(cfg.Script)

	scenario := game.NewScenario(128, 128)
	game.SetupMatch(&scenario)

	maxTicks := cfg.MaxTicks
	if maxTicks <= 0 {
		maxTicks = 0 // unlimited
	}

	for {
		halt, err := exec.ExecTick(&scenario)
		if err != nil {
			fmt.Fprintf(os.Stderr, "script error at tick %d: %v\n", scenario.Tick, err)
			os.Exit(1)
		}
		if halt {
			break
		}

		scenario.Step()

		if maxTicks > 0 && scenario.Tick >= maxTicks {
			fmt.Fprintf(os.Stderr, "reached max ticks (%d)\n", maxTicks)
			break
		}
	}

	if exec.Failed() {
		os.Exit(1)
	}
}
