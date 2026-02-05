# Units

## Rendering

Units are rendered as wide emoji glyphs (2 terminal columns each), matching the tile width. Player ownership is distinguished by color (player-colored background or foreground tint).

### Unit Glyphs

| Unit | Emoji | Notes |
|---|---|---|
| Peasant | 👷 | Worker -- gathers gold/lumber, builds structures |
| Footman | 💂 | Basic melee soldier |

### Future Unit Ideas

| Unit | Emoji | Role |
|---|---|---|
| Mage | 🧙 | Ranged caster |
| Knight | 🏇 | Heavy cavalry |
| Archer/Elf | 🧝 | Ranged physical |
| Ninja/Rogue | 🥷 | Fast stealth unit |
| Prince/Hero | 🤴 | Hero/commander |

### Villain / Undead Faction (if added)

| Unit | Emoji | Role |
|---|---|---|
| Zombie | 🧟 | Basic melee (undead footman) |
| Vampire | 🧛 | Lifesteal melee |
| Supervillain | 🦹 | Hero unit |

## Stats

All values are int32. Speeds and ranges use fixed-point (1000 = 1 tile).

| Unit | HP | Speed | Damage | Range | Attack Speed | Special |
|---|---|---|---|---|---|---|
| Peasant | 30 | 80 | -- | -- | -- | Gathers resources (rate 10, carry 100) |
| Footman | 60 | 100 | 6 | 1500 | 10 ticks | -- |

## Buildings

| Building | Glyph | HP | Role |
|---|---|---|---|
| Town Hall | H | 1200 | Resource drop-off, builds peasants |
| Barracks | B | 800 | Builds footmen |
| Farm | W | 400 | Supply / population |
