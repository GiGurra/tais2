package game

// All values are int32. Positions use fixed-point: 1000 = 1 tile.

type Position struct {
	X, Y int32
}

type Health struct {
	Current, Max int32
}

type Movement struct {
	Speed int32 // tiles/tick in fixed-point (1000 = 1 tile/tick)
}

type Combat struct {
	Damage      int32
	Range       int32 // fixed-point
	AttackSpeed int32 // ticks between attacks
	Cooldown    int32 // ticks remaining
}

type Owner struct {
	PlayerID int32
}

type UnitTypeComp struct {
	Type UnitType
}

type BuildQueue struct {
	Queue [5]UnitType
	Len   int32
	Ticks int32 // ticks into current build
}

type Renderable struct {
	Glyph rune
	Color int32
}

type ResourceGatherer struct {
	GatherRate int32
	Carrying   int32
	CarryMax   int32
}

type ResourceStore struct {
	Gold   int32
	Lumber int32
}

type MoveTarget struct {
	X, Y int32 // fixed-point target position
}
