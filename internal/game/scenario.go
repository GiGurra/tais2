package game

// Scenario holds all game state for a single match.
type Scenario struct {
	Terrain Terrain
	World   World
	Tick    int32
}

func NewScenario(mapWidth, mapHeight int32) Scenario {
	return Scenario{
		Terrain: NewTerrain(mapWidth, mapHeight),
		World:   NewWorld(),
	}
}
