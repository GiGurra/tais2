package game

// StartPosition represents a player's starting tile coordinate.
type StartPosition struct {
	X, Y int32
}

// Scenario holds all game state for a single match.
type Scenario struct {
	Terrain        Terrain
	World          World
	Tick           int32
	TickRate       int32 // ticks per second (default 10)
	StartPositions [2]StartPosition
}

func NewScenario(mapWidth, mapHeight int32) Scenario {
	t := NewTerrain(mapWidth, mapHeight)
	sp := generateTerrain(&t)
	return Scenario{
		Terrain:        t,
		World:          NewWorld(),
		TickRate:       10,
		StartPositions: sp,
	}
}

// generateTerrain orchestrates all map generation phases in order.
func generateTerrain(t *Terrain) [2]StartPosition {
	generateForests(t)
	generateWater(t)
	generateGoldMines(t)
	generateDirtPaths(t)
	return generateStartPositions(t)
}

// generateForests creates clustered forests using two-layer deterministic noise.
func generateForests(t *Terrain) {
	for y := int32(0); y < t.Height; y++ {
		for x := int32(0); x < t.Width; x++ {
			// Coarse layer: ~40% of 8x8 zones are forest zones
			if tileHash(x/8, y/8)%100 < 40 {
				// Fine layer: 70% fill within forest zones
				if tileHash(x, y)%100 < 70 {
					t.Set(x, y, Forest)
				}
			}
		}
	}
}

// generateWater creates a central pond with noise-modulated edges.
func generateWater(t *Terrain) {
	cx := t.Width / 2
	cy := t.Height / 2
	radius := int32(12)
	radiusSq := radius * radius

	for y := cy - radius - 2; y <= cy+radius+2; y++ {
		for x := cx - radius - 2; x <= cx+radius+2; x++ {
			if !t.InBounds(x, y) {
				continue
			}
			dx := x - cx
			dy := y - cy
			distSq := dx*dx + dy*dy
			noise := tileHash(x*3, y*3) % 40 * 4
			if distSq < radiusSq-noise {
				t.Set(x, y, Water)
			}
		}
	}
}

// generateGoldMines places 6 gold mines at symmetric positions.
func generateGoldMines(t *Terrain) {
	mines := [][2]int32{
		{20, 20},   // near P1
		{40, 25},   // near P1
		{55, 55},   // contested center
		{73, 73},   // contested center
		{88, 103},  // near P2
		{108, 108}, // near P2
	}

	for _, m := range mines {
		mx, my := m[0], m[1]

		// Clear a 4x4 dirt border around the mine for visibility
		for dy := int32(-1); dy <= 2; dy++ {
			for dx := int32(-1); dx <= 2; dx++ {
				px, py := mx+dx, my+dy
				if t.InBounds(px, py) {
					t.Set(px, py, Dirt)
				}
			}
		}

		// Place 2x2 gold mine patch
		for dy := int32(0); dy <= 1; dy++ {
			for dx := int32(0); dx <= 1; dx++ {
				px, py := mx+dx, my+dy
				if t.InBounds(px, py) {
					t.Set(px, py, GoldMine)
				}
			}
		}
	}
}

// generateDirtPaths draws 2-tile-wide dirt paths from each start area toward center mines.
func generateDirtPaths(t *Terrain) {
	// P1 (15,15) → contested mine (55,55)
	drawDirtPath(t, 15, 15, 55, 55)
	// P2 (width-16, height-16) → contested mine (73,73)
	drawDirtPath(t, t.Width-16, t.Height-16, 73, 73)
}

// drawDirtPath draws a 2-tile-wide dirt line between two points.
// Only overwrites Grass or Forest tiles.
func drawDirtPath(t *Terrain, x0, y0, x1, y1 int32) {
	dx := abs32(x1 - x0)
	dy := abs32(y1 - y0)
	sx := int32(1)
	if x0 > x1 {
		sx = -1
	}
	sy := int32(1)
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy
	x, y := x0, y0

	for {
		// Place 2-tile-wide path
		for oy := int32(0); oy <= 1; oy++ {
			for ox := int32(0); ox <= 1; ox++ {
				px, py := x+ox, y+oy
				if t.InBounds(px, py) {
					tt := t.At(px, py)
					if tt == Grass || tt == Forest {
						t.Set(px, py, Dirt)
					}
				}
			}
		}

		if x == x1 && y == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

// generateStartPositions clears radius-10 circles to grass at each start position.
func generateStartPositions(t *Terrain) [2]StartPosition {
	p1 := StartPosition{X: 15, Y: 15}
	p2 := StartPosition{X: t.Width - 16, Y: t.Height - 16}
	radius := int32(10)
	radiusSq := radius * radius

	for _, sp := range []StartPosition{p1, p2} {
		for y := sp.Y - radius; y <= sp.Y+radius; y++ {
			for x := sp.X - radius; x <= sp.X+radius; x++ {
				if !t.InBounds(x, y) {
					continue
				}
				dx := x - sp.X
				dy := y - sp.Y
				if dx*dx+dy*dy <= radiusSq {
					t.Set(x, y, Grass)
				}
			}
		}
	}

	return [2]StartPosition{p1, p2}
}

// abs32 returns the absolute value of an int32.
func abs32(v int32) int32 {
	if v < 0 {
		return -v
	}
	return v
}

// Step advances the simulation by one tick.
func (s *Scenario) Step() {
	s.Tick++
	// Future systems called here in fixed order
}

// SetupMatch spawns starting units for both players.
// For each player: 1 TownHall at the start position, 4 Peasants offset by 2 tiles.
func SetupMatch(s *Scenario) {
	for playerID := int32(0); playerID < 2; playerID++ {
		sp := s.StartPositions[playerID]
		cx := sp.X * 1000 // convert to fixed-point
		cy := sp.Y * 1000

		s.World.Spawn(UnitTownHall, cx, cy, playerID)

		// 4 Peasants at cardinal offsets of 2 tiles
		offsets := [4][2]int32{
			{0, -2000}, // north
			{0, 2000},  // south
			{-2000, 0}, // west
			{2000, 0},  // east
		}
		for _, off := range offsets {
			s.World.Spawn(UnitPeasant, cx+off[0], cy+off[1], playerID)
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
