# Terrain & Map Generation

## Terrain Types

| Type | ID | Viewport | Minimap | Description |
|---|---|---|---|---|
| Grass | 0 | `..` green | ` ` | Default walkable ground |
| Water | 1 | `~~` blue | `~` | Impassable (pond, rivers) |
| Forest | 2 | 🌲 green | `T` | Harvestable lumber, blocks movement |
| Dirt | 3 | `..` brown | `.` | Walkable path / cleared ground |
| Mountain | 4 | `^^` gray | `^` | Impassable high ground |
| Gold Mine | 5 | `$$` yellow/brown | `$` | Harvestable gold resource |

## Map Generation Phases

Maps are generated procedurally in 5 phases (later phases override earlier):

1. **Clustered forests** -- two-layer deterministic noise (coarse 8x8 zones at 40% + fine per-tile at 70%) producing ~28% natural-looking forest coverage
2. **Central pond** -- circular water body at map center (radius ~12) with noise-modulated edges for organic shape
3. **Gold mines** -- 6 mines at symmetric positions (2 near each player, 2 contested in center). Each is a 2x2 gold patch with a dirt border
4. **Dirt paths** -- 2-tile-wide Bresenham lines from each start area toward the contested center mines
5. **Starting positions** -- radius-10 circles cleared to grass, ensuring each player has open space

All generation is fully deterministic using an integer hash function (no RNG, no floats).

## Map Layout (128x128)

```
P1 (15,15)          Gold (40,25)
    *----path----------->
  Gold (20,20)              Pond (64,64)
                        Gold (55,55)
                        Gold (73,73)
              <-----------path----*
        Gold (88,103)         P2 (112,112)
                           Gold (108,108)
```
