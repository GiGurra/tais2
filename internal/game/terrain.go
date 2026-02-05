package game

// TerrainType represents the type of terrain on a map tile.
type TerrainType byte

const (
	Grass    TerrainType = 0
	Water    TerrainType = 1
	Forest   TerrainType = 2
	Dirt     TerrainType = 3
	Mountain TerrainType = 4
	GoldMine TerrainType = 5
)

// Terrain is a 2D grid of terrain tiles stored in row-major order.
type Terrain struct {
	Width  int32
	Height int32
	Tiles  []byte // index = y*Width + x
}

func NewTerrain(w, h int32) Terrain {
	return Terrain{
		Width:  w,
		Height: h,
		Tiles:  make([]byte, w*h),
	}
}

func (t *Terrain) InBounds(x, y int32) bool {
	return x >= 0 && x < t.Width && y >= 0 && y < t.Height
}

func (t *Terrain) At(x, y int32) TerrainType {
	return TerrainType(t.Tiles[y*t.Width+x])
}

func (t *Terrain) Set(x, y int32, tt TerrainType) {
	t.Tiles[y*t.Width+x] = byte(tt)
}

// IsWalkable returns true if the tile at (x, y) can be walked on.
// Walkable terrain: Grass, Dirt. Out-of-bounds is not walkable.
func (t *Terrain) IsWalkable(x, y int32) bool {
	if !t.InBounds(x, y) {
		return false
	}
	tt := t.At(x, y)
	return tt == Grass || tt == Dirt
}
