package script

import (
	"bytes"
	"strings"
	"testing"

	"github.com/GiGurra/tais2/internal/game"
)

func TestParseBasic(t *testing.T) {
	input := `
# A test script
@0 label peon0 = player 0 peasant 0
@0 label hall0 = player 0 townhall 0
@0 move peon0 20 20
@0 snapshot
@50 assert peon0 pos-near 20 20 2
@50 assert hall0 alive
@50 snapshot
@50 halt
`
	s, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(s.Commands) != 8 {
		t.Fatalf("expected 8 commands, got %d", len(s.Commands))
	}

	// Verify sorting: all tick-0 commands before tick-50.
	for i, c := range s.Commands {
		if i > 0 && c.Tick < s.Commands[i-1].Tick {
			t.Errorf("commands not sorted by tick at index %d", i)
		}
	}

	// Check label command.
	c0 := s.Commands[0]
	if c0.Type != CmdLabel || c0.LabelName != "peon0" {
		t.Errorf("expected label peon0, got type=%d name=%q", c0.Type, c0.LabelName)
	}
	if c0.Query.PlayerID != 0 || c0.Query.UnitType != game.UnitPeasant || c0.Query.Index != 0 {
		t.Errorf("unexpected query: %+v", c0.Query)
	}

	// Check move command.
	c2 := s.Commands[2]
	if c2.Type != CmdMove || c2.LabelName != "peon0" || c2.TileX != 20 || c2.TileY != 20 {
		t.Errorf("unexpected move: %+v", c2)
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no tick prefix", "label foo = player 0 peasant 0"},
		{"bad tick", "@abc label foo = player 0 peasant 0"},
		{"unknown command", "@0 explode"},
		{"bad label syntax", "@0 label foo"},
		{"bad unit type", "@0 label foo = player 0 dragon 0"},
		{"bad assert kind", "@0 assert foo flying"},
		{"pos-near missing args", "@0 assert foo pos-near 10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(strings.NewReader(tt.input))
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func setupScenario() *game.Scenario {
	s := game.NewScenario(128, 128)
	game.SetupMatch(&s)
	return &s
}

func TestExecutorLabelAndSnapshot(t *testing.T) {
	input := `
@0 label peon0 = player 0 peasant 0
@0 label hall0 = player 0 townhall 0
@0 snapshot
`
	sc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scenario := setupScenario()
	var buf bytes.Buffer
	exec := NewExecutor(sc, &buf)

	halt, err := exec.ExecTick(scenario)
	if err != nil {
		t.Fatalf("ExecTick error: %v", err)
	}
	if halt {
		t.Error("unexpected halt")
	}

	out := buf.String()
	if !strings.Contains(out, "SNAPSHOT tick=0") {
		t.Errorf("snapshot output missing header, got:\n%s", out)
	}
	if !strings.Contains(out, "=== END ===") {
		t.Errorf("snapshot output missing footer, got:\n%s", out)
	}
}

func TestExecutorMoveAndAssert(t *testing.T) {
	// Move a peasant to tile 20,20 and assert after enough ticks.
	input := `
@0 label peon0 = player 0 peasant 0
@0 move peon0 20 20
@500 assert peon0 pos-near 20 20 2
@500 halt
`
	sc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scenario := setupScenario()
	var buf bytes.Buffer
	exec := NewExecutor(sc, &buf)

	for i := range int32(501) {
		halt, err := exec.ExecTick(scenario)
		if err != nil {
			t.Fatalf("tick %d: ExecTick error: %v", i, err)
		}
		if halt {
			if i != 500 {
				t.Errorf("halt at tick %d, expected 500", i)
			}
			break
		}
		scenario.Step()
	}

	if exec.Failed() {
		t.Errorf("assertion failed:\n%s", buf.String())
	}
}

func TestExecutorAssertAlive(t *testing.T) {
	input := `
@0 label hall0 = player 0 townhall 0
@0 assert hall0 alive
@0 halt
`
	sc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scenario := setupScenario()
	var buf bytes.Buffer
	exec := NewExecutor(sc, &buf)

	halt, err := exec.ExecTick(scenario)
	if err != nil {
		t.Fatalf("ExecTick error: %v", err)
	}
	if !halt {
		t.Error("expected halt")
	}
	if exec.Failed() {
		t.Errorf("assertion failed:\n%s", buf.String())
	}
}

func TestExecutorHalt(t *testing.T) {
	input := `@0 halt`
	sc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scenario := setupScenario()
	var buf bytes.Buffer
	exec := NewExecutor(sc, &buf)

	halt, err := exec.ExecTick(scenario)
	if err != nil {
		t.Fatalf("ExecTick error: %v", err)
	}
	if !halt {
		t.Fatal("expected halt=true")
	}
}

func TestExecutorLabelNotFound(t *testing.T) {
	// Requesting peasant index 99 which doesn't exist.
	input := `@0 label p99 = player 0 peasant 99`
	sc, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	scenario := setupScenario()
	var buf bytes.Buffer
	exec := NewExecutor(sc, &buf)

	_, err = exec.ExecTick(scenario)
	if err == nil {
		t.Fatal("expected error for missing entity")
	}
}
