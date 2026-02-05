package game

type UnitType int32

const (
	UnitNone     UnitType = 0
	UnitPeasant  UnitType = 1
	UnitFootman  UnitType = 2
	UnitTownHall UnitType = 3
	UnitBarracks UnitType = 4
	UnitFarm     UnitType = 5
)

// Archetype holds default component values and a mask indicating which
// components are present for a given UnitType.
type Archetype struct {
	Mask             ComponentMask
	Health           Health
	Movement         Movement
	Combat           Combat
	Renderable       Renderable
	ResourceGatherer ResourceGatherer
	ResourceStore    ResourceStore
	BuildQueue       BuildQueue
}

var archetypes = map[UnitType]Archetype{
	UnitPeasant: {
		Mask: MaskPosition | MaskHealth | MaskMovement | MaskOwner | MaskUnitType | MaskRenderable | MaskResourceGatherer,
		Health:     Health{Current: 30, Max: 30},
		Movement:   Movement{Speed: 80},
		Renderable: Renderable{Glyph: 'P', Color: 7},
		ResourceGatherer: ResourceGatherer{
			GatherRate: 10,
			CarryMax:   100,
		},
	},
	UnitFootman: {
		Mask: MaskPosition | MaskHealth | MaskMovement | MaskCombat | MaskOwner | MaskUnitType | MaskRenderable,
		Health:     Health{Current: 60, Max: 60},
		Movement:   Movement{Speed: 100},
		Combat:     Combat{Damage: 6, Range: 1500, AttackSpeed: 10},
		Renderable: Renderable{Glyph: 'F', Color: 1},
	},
	UnitTownHall: {
		Mask: MaskPosition | MaskHealth | MaskOwner | MaskUnitType | MaskRenderable | MaskResourceStore | MaskBuildQueue,
		Health:        Health{Current: 1200, Max: 1200},
		Renderable:    Renderable{Glyph: 'H', Color: 3},
		ResourceStore: ResourceStore{Gold: 0, Lumber: 0},
	},
	UnitBarracks: {
		Mask: MaskPosition | MaskHealth | MaskOwner | MaskUnitType | MaskRenderable | MaskBuildQueue,
		Health:     Health{Current: 800, Max: 800},
		Renderable: Renderable{Glyph: 'B', Color: 3},
	},
	UnitFarm: {
		Mask: MaskPosition | MaskHealth | MaskOwner | MaskUnitType | MaskRenderable,
		Health:     Health{Current: 400, Max: 400},
		Renderable: Renderable{Glyph: 'W', Color: 2},
	},
}

func GetArchetype(ut UnitType) (Archetype, bool) {
	a, ok := archetypes[ut]
	return a, ok
}
