package game

// Scenario holds all game state for a single match.
type Scenario struct {
	Terrain Terrain
	World   World
	Tick    int32
}

func NewScenario(mapWidth, mapHeight int32) Scenario {
	t := NewTerrain(mapWidth, mapHeight)
	scatterTrees(&t)
	return Scenario{
		Terrain: t,
		World:   NewWorld(),
	}
}

// scatterTrees places forest tiles using a simple deterministic hash.
func scatterTrees(t *Terrain) {
	for y := int32(0); y < t.Height; y++ {
		for x := int32(0); x < t.Width; x++ {
			h := tileHash(x, y)
			if h%5 == 0 { // ~20% forest coverage
				t.Set(x, y, Forest)
			}
		}
	}
}

// tileHash returns a deterministic pseudo-random value for a coordinate.
func tileHash(x, y int32) int32 {
	// Simple integer hash (no floats, fully deterministic)
	h := x*374761393 + y*668265263
	h = (h ^ (h >> 13)) * 1274126177
	h = h ^ (h >> 16)
	if h < 0 {
		h = -h
	}
	return h
}
