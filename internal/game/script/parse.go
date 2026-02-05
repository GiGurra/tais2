package script

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/GiGurra/tais2/internal/game"
)

// Parse reads a script from r and returns a sorted Script.
func Parse(r io.Reader) (*Script, error) {
	var cmds []ScriptCommand
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Line must start with @TICK.
		if !strings.HasPrefix(line, "@") {
			return nil, fmt.Errorf("line %d: expected @tick prefix, got %q", lineNum, line)
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: incomplete command", lineNum)
		}

		tick, err := strconv.ParseInt(fields[0][1:], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid tick %q: %w", lineNum, fields[0], err)
		}

		keyword := fields[1]
		args := fields[2:]

		cmd := ScriptCommand{
			Tick: int32(tick),
			Line: lineNum,
		}

		switch keyword {
		case "label":
			if err := parseLabel(&cmd, args, lineNum); err != nil {
				return nil, err
			}
		case "move":
			if err := parseMove(&cmd, args, lineNum); err != nil {
				return nil, err
			}
		case "assert":
			if err := parseAssert(&cmd, args, lineNum); err != nil {
				return nil, err
			}
		case "snapshot":
			cmd.Type = CmdSnapshot
		case "halt":
			cmd.Type = CmdHalt
		default:
			return nil, fmt.Errorf("line %d: unknown command %q", lineNum, keyword)
		}

		cmds = append(cmds, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading script: %w", err)
	}

	sort.Slice(cmds, func(i, j int) bool {
		if cmds[i].Tick != cmds[j].Tick {
			return cmds[i].Tick < cmds[j].Tick
		}
		return cmds[i].Line < cmds[j].Line
	})

	return &Script{Commands: cmds}, nil
}

// parseLabel parses: NAME = player PID TYPE INDEX
func parseLabel(cmd *ScriptCommand, args []string, line int) error {
	// label NAME = player PID TYPE INDEX
	if len(args) < 6 || args[1] != "=" || args[2] != "player" {
		return fmt.Errorf("line %d: label syntax: label NAME = player PID TYPE INDEX", line)
	}
	cmd.Type = CmdLabel
	cmd.LabelName = args[0]

	pid, err := strconv.ParseInt(args[3], 10, 32)
	if err != nil {
		return fmt.Errorf("line %d: invalid player ID %q: %w", line, args[3], err)
	}

	ut, err := parseUnitType(args[4])
	if err != nil {
		return fmt.Errorf("line %d: %w", line, err)
	}

	idx, err := strconv.ParseInt(args[5], 10, 32)
	if err != nil {
		return fmt.Errorf("line %d: invalid index %q: %w", line, args[5], err)
	}

	cmd.Query = EntityQuery{
		PlayerID: int32(pid),
		UnitType: ut,
		Index:    int32(idx),
	}
	return nil
}

// parseMove parses: LABEL TX TY
func parseMove(cmd *ScriptCommand, args []string, line int) error {
	if len(args) < 3 {
		return fmt.Errorf("line %d: move syntax: move LABEL TX TY", line)
	}
	cmd.Type = CmdMove
	cmd.LabelName = args[0]

	tx, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return fmt.Errorf("line %d: invalid tile X %q: %w", line, args[1], err)
	}
	ty, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		return fmt.Errorf("line %d: invalid tile Y %q: %w", line, args[2], err)
	}

	cmd.TileX = int32(tx)
	cmd.TileY = int32(ty)
	return nil
}

// parseAssert parses: LABEL alive|dead|pos-near [X Y TOL]
func parseAssert(cmd *ScriptCommand, args []string, line int) error {
	if len(args) < 2 {
		return fmt.Errorf("line %d: assert syntax: assert LABEL alive|dead|pos-near [X Y TOL]", line)
	}
	cmd.Type = CmdAssert
	cmd.LabelName = args[0]
	cmd.AssertKind = args[1]

	switch cmd.AssertKind {
	case "alive", "dead":
		// No extra args needed.
	case "pos-near":
		if len(args) < 5 {
			return fmt.Errorf("line %d: assert pos-near syntax: assert LABEL pos-near X Y TOL", line)
		}
		x, err := strconv.ParseInt(args[2], 10, 32)
		if err != nil {
			return fmt.Errorf("line %d: invalid X %q: %w", line, args[2], err)
		}
		y, err := strconv.ParseInt(args[3], 10, 32)
		if err != nil {
			return fmt.Errorf("line %d: invalid Y %q: %w", line, args[3], err)
		}
		tol, err := strconv.ParseInt(args[4], 10, 32)
		if err != nil {
			return fmt.Errorf("line %d: invalid tolerance %q: %w", line, args[4], err)
		}
		cmd.TileX = int32(x)
		cmd.TileY = int32(y)
		cmd.Tolerance = int32(tol)
	default:
		return fmt.Errorf("line %d: unknown assert kind %q (want alive, dead, or pos-near)", line, cmd.AssertKind)
	}

	return nil
}

// parseUnitType maps human-readable names to game.UnitType.
func parseUnitType(s string) (game.UnitType, error) {
	switch strings.ToLower(s) {
	case "peasant":
		return game.UnitPeasant, nil
	case "footman":
		return game.UnitFootman, nil
	case "townhall":
		return game.UnitTownHall, nil
	case "barracks":
		return game.UnitBarracks, nil
	case "farm":
		return game.UnitFarm, nil
	default:
		return game.UnitNone, fmt.Errorf("unknown unit type %q", s)
	}
}
