# CLAUDE.md

This file helps AI assistants and humans work on the platformer codebase.

## Overview

2D platformer in **Go** using **Ebitengine** (v2.9.0). Single-player, no enemies, 3 levels. Goal: reach the gold goal zone at the end of each level to advance. Beat all 3 to win.

- **Controls**: W = jump, A = left, D = right, Tab = cycle player shape
- **Screen**: 1280x720 (was 640x360, doubled in session 1)
- **Player shapes**: Purple circle, red triangle, orange hexagon -- all with eyes, all roll visually
- **Levels**: Each is 4 screens wide. Level 1 has full ground + platforms. Level 2 has floating platforms. Level 3 is precision platforming (small platforms, big gaps).

## Commands

```bash
go mod tidy          # resolve dependencies
go run ./cmd/game    # run the game
go build -o platformer ./cmd/game  # compile binary
```

## Architecture

```
cmd/game/main.go           Entry point: creates Game, sets window 1280x720, runs ebiten.RunGame
internal/game/game.go      Game struct (ebiten.Game): Update, Draw, Layout. Holds player, level, camera, state, levelNum. Chains levels 1→2→3.
internal/player/player.go  Player struct: position, velocity, rotation, shape. Input (W/A/D/Tab), gravity, coyote time, jump buffer. Pre-rendered shape images (circle, triangle, hexagon) with vector drawing.
internal/level/level.go    Level struct: platforms ([]image.Rectangle), goal rect, start pos, deathY. FirstLevel(), SecondLevel(), ThirdLevel(). ResolveCollision() with multi-pass AABB. InGoal() for win check.
internal/camera/camera.go  Camera struct: X,Y world position. Smooth lerp follow + clamped to level bounds. WorldToScreen() for drawing.
```

## Key design decisions

- **Y-down coordinate system** (Ebitengine standard). Smaller Y = higher on screen.
- **Fixed dt = 1/60**: Ebitengine runs Update at 60 TPS by default; we use `dt := 1.0/60.0`.
- **Levels are data-defined in Go** (slices of `image.Rectangle`). No external tilemap files.
- **Collision**: Multi-pass (up to 4) AABB overlap with minimum-penetration resolution. Game calls `player.Update(dt)` then `level.ResolveCollision()` then `player.TryJump()`.
- **Player physics**: MoveSpeed=280, JumpVelocity=-420, Gravity=980. Max jump height ~90px. All level jumps must have upward height difference <= ~80px to be achievable.
- **Shape drawing**: Uses `ebiten/v2/vector` package. Circles via `DrawFilledCircle`, polygons (triangle, hexagon) via `vector.Path` + `DrawTriangles` with a 1x1 white pixel source image.

## Git branches

- **main**: Primary branch with all core features
- **feature/missionary-megans-updates**: Feature branch with megan's additions - sunset background, Escape (pina colada song) lyrics, and decorative palm trees

## Screen size and level dimensions

- Screen: 1280x720
- Level width: `screenW * 4` = 5120
- Level height: `screenH * 2` = 1440
- Floor Y: `levelHeight - 48` = 1392

## Conventions

- All game code under `cmd/` and `internal/`
- Units in pixels
- Colors use `color.RGBA{}` (not hex int literals -- Ebitengine's `Fill()` requires `color.Color`)
- Player bounding box is 28x28 (diameter of Radius=14 circle)
