# Networking: Lockstep Simulation

## The standard approach: Turn-based lockstep

Used by Age of Empires, StarCraft: Brood War, and most classic RTS games.

### Core idea

- Time is divided into **command turns**, which run at a lower frequency than simulation ticks. E.g. a command turn every 200ms (5/sec), while simulation runs at 10 ticks/sec — so each command turn spans 2 ticks.
- Each player sends their command list (possibly empty) for turn N **ahead of time**, before the simulation reaches turn N.
- The simulation only advances into turn N's ticks once **all players' commands for turn N have arrived**.
- If a player's commands haven't arrived yet, everyone waits (the classic RTS "lag" stutter).

The **input delay** is the key tuning knob: if a player issues a command at turn 5, it gets scheduled for turn 5+D where D is the agreed-upon delay in command turns. This gives the network time to deliver it to all peers before the simulation needs it.

### Why this works for RTS

Because the simulation is deterministic (same inputs -> same outputs), we only need to send **commands** over the network, not full game state. This keeps bandwidth tiny even with thousands of units.

Prerequisites (already met by tais2):
- Deterministic simulation (no floats, no random without shared seed)
- All game logic uses int32
- Fixed-point math for fractional values

## Design for tais2

### New fields on Scenario

```go
type Scenario struct {
    // ... existing fields ...
    TickRate        int32  // sim ticks/sec (already exists, default 10)
    CommandInterval int32  // ticks per command turn (e.g. 2)
    InputDelay      int32  // command turns ahead to schedule (e.g. 2)
}
```

### CommandTurn message

```go
type Command struct {
    Type     CommandType
    EntityID EntityID
    TargetX  int32  // fixed-point
    TargetY  int32  // fixed-point
    // ... other fields as needed
}

type CommandTurn struct {
    Turn     int32     // which command turn this is for
    PlayerID int32
    Commands []Command // can be empty — acts as heartbeat
}
```

Every player must send one CommandTurn per turn, even if empty. That's the "heartbeat" that lets everyone know it's safe to advance.

### Simulation loop

```go
func (s *Scenario) Step() {
    if s.Tick % s.CommandInterval == 0 {
        s.applyCommands(s.Tick / s.CommandInterval)
    }
    // run systems
    s.Tick++
}
```

Commands are applied *before* systems run for that tick. This keeps the "apply then simulate" ordering deterministic and replayable.

### Singleplayer vs multiplayer

- **Singleplayer**: InputDelay=0, commands go straight into the buffer. No waiting.
- **Multiplayer**: InputDelay=2+, commands are scheduled ahead. Simulation gates on receiving all players' CommandTurns before advancing past a command turn boundary.

The UI layer handles the difference — the simulation itself doesn't know or care about networking.

## Adaptive delay

Rather than a fixed InputDelay, dynamically adjust based on observed round-trip times:

- Start with D=2 command turns
- If commands arrive consistently early, decrease D (snappier feel)
- If stalls occur (waiting on a player), increase D
- Clamp to some min/max (e.g. 1-8 command turns)

Age of Empires 2 did this — the "speed" settings were really just changing the command turn duration and input delay.

## Replays come free

Since the simulation is deterministic, recording the stream of CommandTurn messages is a complete replay. To replay: create the same scenario, feed the same commands at the same turns, get the same result.

## Concrete example with default settings

```
TickRate:        10   (100ms per tick)
CommandInterval:  2   (command turn every 2 ticks = 200ms)
InputDelay:       2   (commands scheduled 2 command turns ahead = 400ms)

Timeline:
  Tick 0 (CT 0): apply commands for CT 0, simulate
  Tick 1:        simulate
  Tick 2 (CT 1): apply commands for CT 1, simulate
  Tick 3:        simulate
  ...

Player issues "move unit" at CT 3 -> scheduled for CT 5
All players must receive CT 5 commands before tick 10
```

At 5 command turns/sec with 2-turn delay, perceived input latency is ~400ms — typical for RTS and barely noticeable for strategic commands.
