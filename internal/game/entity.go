package game

const MaxEntities = 4096

// ComponentMask is a bitmask indicating which components an entity has.
type ComponentMask uint16

const (
	MaskPosition         ComponentMask = 1 << iota // 0
	MaskHealth                                     // 1
	MaskMovement                                   // 2
	MaskCombat                                     // 3
	MaskOwner                                      // 4
	MaskUnitType                                   // 5
	MaskBuildQueue                                 // 6
	MaskRenderable                                 // 7
	MaskResourceGatherer                           // 8
	MaskResourceStore                              // 9
	MaskMoveTarget                                 // 10
)

// EntityID packs a generation (high 16 bits) and a slot index (low 16 bits)
// into a single int32.
type EntityID int32

func MakeEntityID(index, gen int32) EntityID {
	return EntityID((gen << 16) | (index & 0xFFFF))
}

func (id EntityID) Index() int32 { return int32(id) & 0xFFFF }
func (id EntityID) Gen() int32   { return int32(id) >> 16 }

// Entity holds per-slot metadata.
type Entity struct {
	Alive bool
	Mask  ComponentMask
	Gen   int32
}

// World is a struct-of-arrays ECS container.
type World struct {
	Entities [MaxEntities]Entity

	// Component arrays — parallel to Entities.
	Position         [MaxEntities]Position
	Health           [MaxEntities]Health
	Movement         [MaxEntities]Movement
	Combat           [MaxEntities]Combat
	Owner            [MaxEntities]Owner
	UnitTypeComp     [MaxEntities]UnitTypeComp
	BuildQueue       [MaxEntities]BuildQueue
	Renderable       [MaxEntities]Renderable
	ResourceGatherer [MaxEntities]ResourceGatherer
	ResourceStore    [MaxEntities]ResourceStore
	MoveTarget       [MaxEntities]MoveTarget

	freeList  [MaxEntities]int32
	freeCount int32

	AliveCount int32
}

func NewWorld() World {
	var w World
	// Fill free list so slot 0 is popped first (stack order).
	for i := int32(MaxEntities - 1); i >= 0; i-- {
		w.freeList[w.freeCount] = i
		w.freeCount++
	}
	return w
}

// Spawn creates a new entity from the given archetype at the specified position.
func (w *World) Spawn(ut UnitType, x, y, playerID int32) EntityID {
	if w.freeCount == 0 {
		return -1
	}

	arch, ok := GetArchetype(ut)
	if !ok {
		return -1
	}

	// Pop a free slot.
	w.freeCount--
	idx := w.freeList[w.freeCount]

	// Zero all component slots to prevent stale data.
	w.Position[idx] = Position{}
	w.Health[idx] = Health{}
	w.Movement[idx] = Movement{}
	w.Combat[idx] = Combat{}
	w.Owner[idx] = Owner{}
	w.UnitTypeComp[idx] = UnitTypeComp{}
	w.BuildQueue[idx] = BuildQueue{}
	w.Renderable[idx] = Renderable{}
	w.ResourceGatherer[idx] = ResourceGatherer{}
	w.ResourceStore[idx] = ResourceStore{}
	w.MoveTarget[idx] = MoveTarget{}

	// Stamp archetype values.
	e := &w.Entities[idx]
	e.Alive = true
	e.Mask = arch.Mask

	w.Position[idx] = Position{X: x, Y: y}
	w.Health[idx] = arch.Health
	w.Movement[idx] = arch.Movement
	w.Combat[idx] = arch.Combat
	w.Owner[idx] = Owner{PlayerID: playerID}
	w.UnitTypeComp[idx] = UnitTypeComp{Type: ut}
	w.Renderable[idx] = arch.Renderable
	w.ResourceGatherer[idx] = arch.ResourceGatherer
	w.ResourceStore[idx] = arch.ResourceStore
	w.BuildQueue[idx] = arch.BuildQueue

	w.AliveCount++

	return MakeEntityID(idx, e.Gen)
}

// Despawn removes an entity, validating its generational ID.
func (w *World) Despawn(id EntityID) bool {
	idx := id.Index()
	if idx < 0 || idx >= MaxEntities {
		return false
	}
	e := &w.Entities[idx]
	if !e.Alive || e.Gen != id.Gen() {
		return false
	}

	e.Alive = false
	e.Mask = 0
	e.Gen++

	w.freeList[w.freeCount] = idx
	w.freeCount++
	w.AliveCount--

	return true
}

// IsAlive checks whether the entity with the given ID is still alive.
func (w *World) IsAlive(id EntityID) bool {
	idx := id.Index()
	if idx < 0 || idx >= MaxEntities {
		return false
	}
	e := &w.Entities[idx]
	return e.Alive && e.Gen == id.Gen()
}

// Each iterates over all alive entities whose mask is a superset of the
// given mask. The callback receives the slot index. Return false to stop.
func (w *World) Each(mask ComponentMask, fn func(idx int32) bool) {
	for i := int32(0); i < MaxEntities; i++ {
		e := &w.Entities[i]
		if e.Alive && e.Mask&mask == mask {
			if !fn(i) {
				return
			}
		}
	}
}

// EachPair iterates all unique pairs of alive entities where entity A matches
// maskA and entity B matches maskB. Return false to stop.
func (w *World) EachPair(maskA, maskB ComponentMask, fn func(a, b int32) bool) {
	for i := int32(0); i < MaxEntities; i++ {
		ei := &w.Entities[i]
		if !ei.Alive || ei.Mask&maskA != maskA {
			continue
		}
		for j := i + 1; j < MaxEntities; j++ {
			ej := &w.Entities[j]
			if !ej.Alive || ej.Mask&maskB != maskB {
				continue
			}
			if !fn(i, j) {
				return
			}
		}
	}
}
