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
