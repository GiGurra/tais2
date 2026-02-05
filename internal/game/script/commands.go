package script

import "github.com/GiGurra/tais2/internal/game"

// CommandType identifies the kind of script command.
type CommandType int

const (
	CmdLabel    CommandType = iota
	CmdMove
	CmdAssert
	CmdSnapshot
	CmdHalt
)

// EntityQuery identifies an entity by player, unit type, and spawn-order index.
type EntityQuery struct {
	PlayerID int32
	UnitType game.UnitType
	Index    int32
}

// ScriptCommand is a single parsed instruction from a script file.
type ScriptCommand struct {
	Tick        int32
	Type        CommandType
	LabelName   string      // label, move, assert
	Query       EntityQuery // label
	TileX       int32       // move, assert pos-near
	TileY       int32       // move, assert pos-near
	AssertKind  string      // "alive", "dead", "pos-near"
	Tolerance   int32       // assert pos-near
	Line        int         // source line number for error messages
}

// Script is a sorted list of commands parsed from a script file.
type Script struct {
	Commands []ScriptCommand // sorted by (Tick, Line)
}
