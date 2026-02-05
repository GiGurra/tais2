# Input & Mouse Handling

## Terminal Right-Click Problem

Right-click is **not reliably available** across terminal emulators. macOS terminals
(iTerm2, Terminal.app) intercept right-click to show their own context menus, even
when the application has mouse reporting enabled. This is a terminal emulator UI
design choice, not a protocol limitation — the xterm mouse protocol fully supports
right-click (button 3).

### Terminal Compatibility

| Terminal      | Right-click passthrough? |
|---------------|--------------------------|
| iTerm2        | No — always shows context menu |
| Terminal.app  | No — same problem |
| Alacritty     | Yes |
| Kitty         | Yes |
| WezTerm       | Yes |
| xterm         | Yes |

### Industry Pattern

No popular TUI app (vim, tmux, lazygit, micro) relies on right-click for core
functionality. They are all keyboard-first with mouse as supplementary.

## tais2 Input Design

We support **both** right-click and keyboard command modes, so the game works on
all terminals.

### Mouse

- **Left-click**: Select entity / execute pending command (e.g. move target)
- **Right-click**: Issue move order on selected unit (works on Alacritty, Kitty, WezTerm)

### Keyboard Command Modes

Press a key to enter a command mode, then left-click to execute. Press `Esc` to
cancel the pending command (returns to normal select mode).

| Key | Mode    | Left-click action               |
|-----|---------|----------------------------------|
| m   | Move    | Move selected unit to clicked tile |

Future command modes: `a` (attack-move), `p` (patrol), `b` (build), etc.

### Camera

- `WASD` / Arrow keys: pan camera (1 tile)
- `Shift`/`Alt`/`Ctrl` + Arrow: pan camera (1 screen)
- `Esc`: return to main menu (when no command pending)
