package script

import (
	"fmt"
	"io"

	"github.com/GiGurra/tais2/internal/game"
)

// Executor runs script commands against a live scenario.
type Executor struct {
	script  *Script
	labels  map[string]game.EntityID
	nextCmd int
	output  io.Writer
	failed  bool
}

// NewExecutor creates an Executor that writes snapshot/error output to w.
func NewExecutor(s *Script, w io.Writer) *Executor {
	return &Executor{
		script: s,
		labels: make(map[string]game.EntityID),
		output: w,
	}
}

// Failed returns true if any assertion has failed.
func (e *Executor) Failed() bool {
	return e.failed
}

// ExecTick runs all commands scheduled at the scenario's current tick.
// It should be called before Step(). Returns halt=true if a halt command was executed.
func (e *Executor) ExecTick(s *game.Scenario) (halt bool, err error) {
	for e.nextCmd < len(e.script.Commands) {
		cmd := &e.script.Commands[e.nextCmd]
		if cmd.Tick > s.Tick {
			break
		}
		if cmd.Tick < s.Tick {
			// Skip commands from past ticks (shouldn't happen with sorted input).
			e.nextCmd++
			continue
		}

		h, err := e.execCommand(cmd, s)
		if err != nil {
			return false, err
		}
		e.nextCmd++
		if h {
			return true, nil
		}
	}
	return false, nil
}

func (e *Executor) execCommand(cmd *ScriptCommand, s *game.Scenario) (halt bool, err error) {
	switch cmd.Type {
	case CmdLabel:
		return false, e.execLabel(cmd, s)
	case CmdMove:
		return false, e.execMove(cmd, s)
	case CmdAssert:
		return false, e.execAssert(cmd, s)
	case CmdSnapshot:
		WriteSnapshot(e.output, s)
		return false, nil
	case CmdHalt:
		return true, nil
	default:
		return false, fmt.Errorf("line %d: unknown command type %d", cmd.Line, cmd.Type)
	}
}

func (e *Executor) execLabel(cmd *ScriptCommand, s *game.Scenario) error {
	id, ok := resolveQuery(&s.World, cmd.Query)
	if !ok {
		return fmt.Errorf("line %d: label %q: no matching entity (player=%d type=%d index=%d)",
			cmd.Line, cmd.LabelName, cmd.Query.PlayerID, cmd.Query.UnitType, cmd.Query.Index)
	}
	e.labels[cmd.LabelName] = id
	return nil
}

func (e *Executor) execMove(cmd *ScriptCommand, s *game.Scenario) error {
	id, ok := e.labels[cmd.LabelName]
	if !ok {
		return fmt.Errorf("line %d: move: unknown label %q", cmd.Line, cmd.LabelName)
	}

	idx := id.Index()
	ent := &s.World.Entities[idx]
	if !ent.Alive || ent.Gen != id.Gen() {
		return fmt.Errorf("line %d: move: entity %q (id=%d) is dead or recycled", cmd.Line, cmd.LabelName, id)
	}

	targetX := cmd.TileX*1000 + 500
	targetY := cmd.TileY*1000 + 500

	s.World.MoveTarget[idx] = game.MoveTarget{X: targetX, Y: targetY}
	ent.Mask |= game.MaskMoveTarget

	return nil
}

func (e *Executor) execAssert(cmd *ScriptCommand, s *game.Scenario) error {
	id, ok := e.labels[cmd.LabelName]
	if !ok {
		return fmt.Errorf("line %d: assert: unknown label %q", cmd.Line, cmd.LabelName)
	}

	idx := id.Index()
	ent := &s.World.Entities[idx]
	alive := ent.Alive && ent.Gen == id.Gen()

	switch cmd.AssertKind {
	case "alive":
		if !alive {
			e.failed = true
			fmt.Fprintf(e.output, "ASSERT FAILED line %d: %q expected alive, but dead\n", cmd.Line, cmd.LabelName)
		}
	case "dead":
		if alive {
			e.failed = true
			fmt.Fprintf(e.output, "ASSERT FAILED line %d: %q expected dead, but alive\n", cmd.Line, cmd.LabelName)
		}
	case "pos-near":
		if !alive {
			e.failed = true
			fmt.Fprintf(e.output, "ASSERT FAILED line %d: %q expected pos-near but entity is dead\n", cmd.Line, cmd.LabelName)
			return nil
		}
		pos := s.World.Position[idx]
		tileX := pos.X / 1000
		tileY := pos.Y / 1000
		dx := tileX - cmd.TileX
		dy := tileY - cmd.TileY
		if dx < 0 {
			dx = -dx
		}
		if dy < 0 {
			dy = -dy
		}
		if dx > cmd.Tolerance || dy > cmd.Tolerance {
			e.failed = true
			fmt.Fprintf(e.output, "ASSERT FAILED line %d: %q expected pos-near %d,%d (tol=%d) but at %d,%d\n",
				cmd.Line, cmd.LabelName, cmd.TileX, cmd.TileY, cmd.Tolerance, tileX, tileY)
		}
	}

	return nil
}

// resolveQuery finds the Nth entity matching (playerID, unitType) by slot order.
func resolveQuery(w *game.World, q EntityQuery) (game.EntityID, bool) {
	mask := game.MaskOwner | game.MaskUnitType
	count := int32(0)
	var result game.EntityID
	found := false

	w.Each(mask, func(idx int32) bool {
		if w.Owner[idx].PlayerID != q.PlayerID {
			return true
		}
		if w.UnitTypeComp[idx].Type != q.UnitType {
			return true
		}
		if count == q.Index {
			result = game.MakeEntityID(idx, w.Entities[idx].Gen)
			found = true
			return false
		}
		count++
		return true
	})

	return result, found
}
