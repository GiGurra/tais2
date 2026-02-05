# tais2 - Terminal RTS

A terminal-based real-time strategy game inspired by Warcraft II / Age of Empires, built with Go and Bubbletea.

## Core Design

- **Deterministic simulation** -- all game logic uses int32, no floats. Fixed-point math (1000 = 1 tile) for sub-tile precision. This is a prerequisite for future lockstep multiplayer.
- **ECS architecture** -- struct-of-arrays entity component system with bitmask queries. Max 4096 entities. Generational entity IDs for safe despawn.
- **Terminal UI** -- Bubbletea framework with a custom framebuffer renderer. Each map tile is 2 terminal columns wide. Minimum terminal size 100x40.

## Game Structure

A **Scenario** holds all state for a match: terrain grid, ECS world, tick counter, and starting positions.

The **Terrain** is a 128x128 tile grid with procedural generation: clustered forests, a central pond, gold mines, dirt paths, and cleared starting areas.

The **World** (ECS) manages all entities -- units and buildings -- with components for position, health, movement, combat, ownership, rendering, resource gathering, and building.

## Screens

- **Main Menu** -- entry point
- **Game View** -- main gameplay screen with viewport, minimap, selection panel, command panel, and HUD
