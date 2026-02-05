package game

import "testing"

func TestSpawnAndComponents(t *testing.T) {
	w := NewWorld()

	id := w.Spawn(UnitPeasant, 5000, 3000, 1)
	if id == -1 {
		t.Fatal("expected valid entity ID")
	}
	if !w.IsAlive(id) {
		t.Fatal("entity should be alive")
	}
	if w.AliveCount != 1 {
		t.Fatalf("expected AliveCount=1, got %d", w.AliveCount)
	}

	idx := id.Index()
	if w.Position[idx].X != 5000 || w.Position[idx].Y != 3000 {
		t.Fatalf("position mismatch: got %v", w.Position[idx])
	}
	if w.Health[idx].Current != 30 || w.Health[idx].Max != 30 {
		t.Fatalf("health mismatch: got %v", w.Health[idx])
	}
	if w.Renderable[idx].Glyph != 'P' {
		t.Fatalf("glyph mismatch: got %c", w.Renderable[idx].Glyph)
	}
	if w.Owner[idx].PlayerID != 1 {
		t.Fatalf("owner mismatch: got %d", w.Owner[idx].PlayerID)
	}
	if w.ResourceGatherer[idx].CarryMax != 100 {
		t.Fatalf("resource gatherer mismatch: got %v", w.ResourceGatherer[idx])
	}
}

func TestDespawnAndFreeListReuse(t *testing.T) {
	w := NewWorld()

	id1 := w.Spawn(UnitFootman, 1000, 1000, 1)
	idx1 := id1.Index()
	w.Despawn(id1)

	if w.IsAlive(id1) {
		t.Fatal("entity should be dead after despawn")
	}
	if w.AliveCount != 0 {
		t.Fatalf("expected AliveCount=0, got %d", w.AliveCount)
	}

	// Next spawn should reuse the same slot.
	id2 := w.Spawn(UnitPeasant, 2000, 2000, 2)
	if id2.Index() != idx1 {
		t.Fatalf("expected slot reuse: got index %d, want %d", id2.Index(), idx1)
	}
	// Generation should have incremented.
	if id2.Gen() <= id1.Gen() {
		t.Fatalf("expected generation increment: got %d, previous %d", id2.Gen(), id1.Gen())
	}
}

func TestStaleEntityID(t *testing.T) {
	w := NewWorld()

	id := w.Spawn(UnitFootman, 0, 0, 1)
	w.Despawn(id)

	// Stale ID should not be alive.
	if w.IsAlive(id) {
		t.Fatal("stale entity ID should not be alive")
	}

	// Stale ID should not be despawnable again.
	if w.Despawn(id) {
		t.Fatal("stale entity ID should not be despawnable")
	}
}

func TestEachByMask(t *testing.T) {
	w := NewWorld()

	w.Spawn(UnitPeasant, 0, 0, 1)  // has Movement + ResourceGatherer
	w.Spawn(UnitFootman, 0, 0, 1)  // has Movement + Combat
	w.Spawn(UnitTownHall, 0, 0, 1) // no Movement

	// Query entities with Movement component.
	count := int32(0)
	w.Each(MaskMovement, func(idx int32) bool {
		count++
		return true
	})
	if count != 2 {
		t.Fatalf("expected 2 entities with Movement, got %d", count)
	}

	// Query entities with Combat component.
	count = 0
	w.Each(MaskCombat, func(idx int32) bool {
		count++
		return true
	})
	if count != 1 {
		t.Fatalf("expected 1 entity with Combat, got %d", count)
	}

	// Query entities with ResourceStore component.
	count = 0
	w.Each(MaskResourceStore, func(idx int32) bool {
		count++
		return true
	})
	if count != 1 {
		t.Fatalf("expected 1 entity with ResourceStore, got %d", count)
	}
}

func TestTerrainAtSetInBounds(t *testing.T) {
	ter := NewTerrain(10, 10)

	if !ter.InBounds(0, 0) {
		t.Fatal("(0,0) should be in bounds")
	}
	if !ter.InBounds(9, 9) {
		t.Fatal("(9,9) should be in bounds")
	}
	if ter.InBounds(10, 0) {
		t.Fatal("(10,0) should be out of bounds")
	}
	if ter.InBounds(-1, 0) {
		t.Fatal("(-1,0) should be out of bounds")
	}

	ter.Set(3, 4, Water)
	if ter.At(3, 4) != Water {
		t.Fatalf("expected Water at (3,4), got %d", ter.At(3, 4))
	}

	ter.Set(0, 0, Mountain)
	if ter.At(0, 0) != Mountain {
		t.Fatalf("expected Mountain at (0,0), got %d", ter.At(0, 0))
	}

	// Default should be Grass (0).
	if ter.At(5, 5) != Grass {
		t.Fatalf("expected Grass at (5,5), got %d", ter.At(5, 5))
	}
}

func TestEntityIDPacking(t *testing.T) {
	id := MakeEntityID(42, 7)
	if id.Index() != 42 {
		t.Fatalf("expected index 42, got %d", id.Index())
	}
	if id.Gen() != 7 {
		t.Fatalf("expected gen 7, got %d", id.Gen())
	}
}

func TestNewScenario(t *testing.T) {
	s := NewScenario(64, 48)
	if s.Terrain.Width != 64 || s.Terrain.Height != 48 {
		t.Fatalf("terrain size mismatch: %dx%d", s.Terrain.Width, s.Terrain.Height)
	}
	if s.Tick != 0 {
		t.Fatalf("expected tick 0, got %d", s.Tick)
	}
	if s.World.AliveCount != 0 {
		t.Fatalf("expected 0 alive, got %d", s.World.AliveCount)
	}
}

func TestTickRate(t *testing.T) {
	s := NewScenario(64, 48)
	if s.TickRate != 10 {
		t.Fatalf("expected default TickRate=10, got %d", s.TickRate)
	}
}

func TestStep(t *testing.T) {
	s := NewScenario(64, 48)
	if s.Tick != 0 {
		t.Fatalf("expected tick 0, got %d", s.Tick)
	}
	s.Step()
	if s.Tick != 1 {
		t.Fatalf("expected tick 1 after Step, got %d", s.Tick)
	}
	s.Step()
	s.Step()
	if s.Tick != 3 {
		t.Fatalf("expected tick 3 after 3 Steps, got %d", s.Tick)
	}
}

func TestTerrainWalkability(t *testing.T) {
	ter := NewTerrain(10, 10)
	// Default is Grass — walkable
	if !ter.IsWalkable(0, 0) {
		t.Fatal("grass should be walkable")
	}

	ter.Set(1, 0, Dirt)
	if !ter.IsWalkable(1, 0) {
		t.Fatal("dirt should be walkable")
	}

	ter.Set(2, 0, Water)
	if ter.IsWalkable(2, 0) {
		t.Fatal("water should not be walkable")
	}

	ter.Set(3, 0, Forest)
	if ter.IsWalkable(3, 0) {
		t.Fatal("forest should not be walkable")
	}

	ter.Set(4, 0, Mountain)
	if ter.IsWalkable(4, 0) {
		t.Fatal("mountain should not be walkable")
	}

	ter.Set(5, 0, GoldMine)
	if ter.IsWalkable(5, 0) {
		t.Fatal("gold mine should not be walkable")
	}

	// Out of bounds
	if ter.IsWalkable(-1, 0) {
		t.Fatal("out of bounds should not be walkable")
	}
	if ter.IsWalkable(10, 0) {
		t.Fatal("out of bounds should not be walkable")
	}
}

func TestMovementSystem(t *testing.T) {
	// Build a scenario with clean grass terrain (no generated features)
	s := Scenario{
		Terrain:  NewTerrain(32, 32),
		World:    NewWorld(),
		TickRate: 10,
	}
	// Terrain is all grass by default — walkable

	// Spawn a peasant at (5,5) in fixed-point
	id := s.World.Spawn(UnitPeasant, 5000, 5000, 0)
	idx := id.Index()

	// Issue move order east: target (10,5) center of tile
	s.World.MoveTarget[idx] = MoveTarget{X: 10500, Y: 5500}
	s.World.Entities[idx].Mask |= MaskMoveTarget

	// Speed is 80 per tick, distance is ~5500 units, so need many ticks
	for i := 0; i < 200; i++ {
		s.Step()
		if s.World.Entities[idx].Mask&MaskMoveTarget == 0 {
			break
		}
	}

	// MoveTarget should be cleared
	if s.World.Entities[idx].Mask&MaskMoveTarget != 0 {
		t.Fatal("MaskMoveTarget should be cleared after arrival")
	}

	// Should be at or near target
	pos := s.World.Position[idx]
	if pos.X != 10500 || pos.Y != 5500 {
		t.Fatalf("expected position (10500,5500), got (%d,%d)", pos.X, pos.Y)
	}
}

func TestMovementBlockedByWater(t *testing.T) {
	s := Scenario{
		Terrain:  NewTerrain(32, 32),
		World:    NewWorld(),
		TickRate: 10,
	}

	// Place water at tile (7,5)
	s.Terrain.Set(7, 5, Water)

	// Spawn peasant at tile (5,5) center
	id := s.World.Spawn(UnitPeasant, 5500, 5500, 0)
	idx := id.Index()

	// Issue move order to (9,5) — must pass through water at (7,5)
	s.World.MoveTarget[idx] = MoveTarget{X: 9500, Y: 5500}
	s.World.Entities[idx].Mask |= MaskMoveTarget

	for i := 0; i < 200; i++ {
		s.Step()
		if s.World.Entities[idx].Mask&MaskMoveTarget == 0 {
			break
		}
	}

	// Should have stopped — MaskMoveTarget cleared
	if s.World.Entities[idx].Mask&MaskMoveTarget != 0 {
		t.Fatal("MaskMoveTarget should be cleared when blocked")
	}

	// Should NOT be at the target
	pos := s.World.Position[idx]
	if pos.X == 9500 && pos.Y == 5500 {
		t.Fatal("peasant should not have reached target through water")
	}

	// Should still be on the left side of the water tile
	if pos.X/1000 >= 7 {
		t.Fatalf("peasant should have stopped before water tile, got tile X=%d", pos.X/1000)
	}
}

func TestIsqrt64(t *testing.T) {
	cases := []struct {
		input int64
		want  int64
	}{
		{0, 0},
		{1, 1},
		{4, 2},
		{9, 3},
		{10, 3},
		{15, 3},
		{16, 4},
		{100, 10},
		{1000000, 1000},
		{1000001, 1000},
	}
	for _, tc := range cases {
		got := isqrt64(tc.input)
		if got != tc.want {
			t.Errorf("isqrt64(%d) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestMoveTargetComponentMask(t *testing.T) {
	w := NewWorld()

	id := w.Spawn(UnitPeasant, 5000, 5000, 0)
	idx := id.Index()

	// Peasant should not have MaskMoveTarget by default
	if w.Entities[idx].Mask&MaskMoveTarget != 0 {
		t.Fatal("peasant should not have MaskMoveTarget on spawn")
	}

	// Add move target dynamically
	w.MoveTarget[idx] = MoveTarget{X: 10000, Y: 10000}
	w.Entities[idx].Mask |= MaskMoveTarget

	// Should be queryable via Each
	count := int32(0)
	w.Each(MaskMoveTarget, func(idx int32) bool {
		count++
		return true
	})
	if count != 1 {
		t.Fatalf("expected 1 entity with MaskMoveTarget, got %d", count)
	}

	// Clear it
	w.Entities[idx].Mask &^= MaskMoveTarget
	count = 0
	w.Each(MaskMoveTarget, func(idx int32) bool {
		count++
		return true
	})
	if count != 0 {
		t.Fatalf("expected 0 entities with MaskMoveTarget after clear, got %d", count)
	}
}

func TestSetupMatch(t *testing.T) {
	s := NewScenario(128, 128)
	if s.World.AliveCount != 0 {
		t.Fatalf("expected 0 alive before SetupMatch, got %d", s.World.AliveCount)
	}

	SetupMatch(&s)

	if s.World.AliveCount != 10 {
		t.Fatalf("expected 10 alive after SetupMatch, got %d", s.World.AliveCount)
	}

	// Count TownHalls and Peasants
	townHalls := int32(0)
	peasants := int32(0)
	s.World.Each(MaskUnitType, func(idx int32) bool {
		switch s.World.UnitTypeComp[idx].Type {
		case UnitTownHall:
			townHalls++
		case UnitPeasant:
			peasants++
		}
		return true
	})

	if townHalls != 2 {
		t.Fatalf("expected 2 TownHalls, got %d", townHalls)
	}
	if peasants != 8 {
		t.Fatalf("expected 8 Peasants, got %d", peasants)
	}

	// Verify player ownership: 5 entities per player
	p0Count := int32(0)
	p1Count := int32(0)
	s.World.Each(MaskOwner, func(idx int32) bool {
		switch s.World.Owner[idx].PlayerID {
		case 0:
			p0Count++
		case 1:
			p1Count++
		}
		return true
	})
	if p0Count != 5 {
		t.Fatalf("expected 5 entities for player 0, got %d", p0Count)
	}
	if p1Count != 5 {
		t.Fatalf("expected 5 entities for player 1, got %d", p1Count)
	}
}
