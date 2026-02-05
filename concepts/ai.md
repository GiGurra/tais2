# AI Architecture

Two layers: low-level unit behavior and high-level strategic commander.

## Low-Level: Unit Behavior (FSM)

Each unit has a finite state machine as an ECS component:

```
UnitAI {
    State      // Idle, MoveTo, AttackTarget, AttackMove, Gather, ReturnResource, Build
    TargetID   // entity to interact with
    TargetX, Y // position to move toward
    OrderQueue // small queue of commands (shift-click)
}
```

A `unitAISystem` runs each tick and handles state transitions:

| State | Behavior | Transitions |
|---|---|---|
| Idle | Stand still, scan for enemies in aggro radius | Enemy found -> AttackTarget |
| MoveTo | Step toward TargetX/Y | Arrived -> Idle, Enemy in path + attack-move -> AttackTarget |
| AttackTarget | Move toward target, attack when in range | Target dead -> Idle, Out of range -> chase |
| AttackMove | Move toward point, engage enemies along the way | Enemy in radius -> AttackTarget, Arrived -> Idle |
| Gather | Move to resource -> harvest ticks -> carry -> ReturnResource | Resource depleted -> find nearest |
| ReturnResource | Walk to nearest Town Hall, deposit | Deposited -> Gather (same resource) |
| Build | Move to build site, tick build progress | Complete -> Idle |

All int32, fully deterministic. Aggro scan is a brute-force distance check (fine at 4096 max entities).

## High-Level: Strategic AI (Commander)

Runs every ~60 ticks (not every frame). Priority-based rule system:

```
StrategicAI {
    Phase        // Opening, Expanding, Attacking, Defending
    DesiredArmy  // target unit counts
    AttackScore  // accumulated "confidence" to attack
    ScoutTimer
}
```

### Weighted priority evaluation

| Priority | Condition | Action |
|---|---|---|
| Critical | Under attack, no defenders | Pull workers, rally army |
| High | 0 workers | Build worker |
| High | Supply capped | Build farm |
| Medium | Workers < 6 | Build worker |
| Medium | No barracks | Build barracks |
| Medium | Gold > threshold | Train footman |
| Low | Army > N and idle | Attack-move toward enemy |
| Low | No known enemy positions | Scout |

Difficulty tuning: adjust thresholds and reaction speed (how often the AI re-evaluates).

### Cheating vs fair play

First pass: AI sees the whole map (no fog of war yet). Later, restrict to scouted info only.

## Implementation order

1. UnitAI component + system with Idle / MoveTo / AttackTarget
2. Aggro scan -- idle units auto-attack enemies within ~5 tile radius
3. Gather loop -- peasants cycle between gold mine and town hall
4. Strategic AI layer on top, issuing commands like a player would
