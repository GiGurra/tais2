# Mouse Cursor Appearance

## Problem

Terminal TUI apps default to an I-beam (text selection) cursor, which feels wrong
for an RTS game. We'd prefer a crosshair or pointer.

## OSC 22 Escape Sequence

The `OSC 22` escape sequence can set the mouse cursor shape using CSS cursor names:

```
\x1b]22;crosshair\x1b\\   # set crosshair
\x1b]22;\x1b\\             # reset to default
```

Supported shapes include: `crosshair`, `pointer`, `default`, `move`, `grab`,
`not-allowed`, various resize cursors, etc.

Bubbletea has no built-in API for this — send the escape sequence directly via
`fmt.Print`.

### Terminal Support

| Terminal   | OSC 22 support? |
|------------|-----------------|
| Kitty      | Yes             |
| foot       | Yes             |
| xterm      | Yes             |
| WezTerm    | In progress     |
| Alacritty  | In progress     |
| iTerm2     | No              |

## Options

### Option 1: OSC 22 with graceful degradation

Send `\x1b]22;crosshair\x1b\\` on game start, reset on exit. Terminals that don't
support it silently ignore the sequence. Zero rendering cost.

Could also change cursor contextually (e.g. `crosshair` in move mode, `default`
otherwise).

### Option 2: Software cursor

Hide the OS cursor (`\x1b[?25l`), track mouse position via `tea.MouseMsg`, draw a
custom character (`+`, `⊕`, `⌖`) at the mouse location in `View()`.

Pros: works on all terminals, full control over appearance.
Cons: rendering complexity, can't overlap content cleanly, potential flicker.

### Option 3: Hybrid

Try OSC 22, offer `--software-cursor` flag as fallback.

## Status

Not implemented. Revisit when other terminals catch up on OSC 22 support, or if
the software cursor approach becomes worthwhile.
