# tais2

Terminal-based RTS game built with Go and Bubbletea.

## Architecture Principles

- **Deterministic simulation**: All game logic must be deterministic. This is a prerequisite for future multiplayer (lockstep networking).
- **No floating point**: Use `int32` for all game values (positions, health, damage, timers, etc.). No `float32` or `float64` in game logic.
- **Fixed-point math**: If fractional values are needed, use fixed-point representation (e.g. multiply by 1000).

## Terminal Requirements

- Minimum terminal size: 100x40
- If terminal drops below minimum, overlay a resize prompt instead of rendering the game.

## Tech Stack

- Go
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
