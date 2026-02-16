
# CLAUDE.md

This file helps AI and humans work on the platformer codebase.

## Overview

2D platformer in **Go** using **Ebitengine**. Single level (4 screens wide), no enemies. Goal: reach the end. Controls: **W** jump, **A** left, **D** right. Camera follows the player.

## Commands

```bash
go mod tidy
go run ./cmd/game
go build -o platformer ./cmd/game
```

## Architecture

Entry point is `cmd/game/main.go`, which creates `game.Game` and runs `ebiten.RunGame`. The `internal/game` package owns the game loop and holds the player, level, and camera. Level 1 is data-defined in `internal/level` (platforms and goal as rectangles); no external tilemap for v1.

- **cmd/game/main.go** — Entry point, window size/title, runs the game.
- **internal/game/game.go** — Implements `ebiten.Game`: Update (input, player, collision, camera, win/death), Draw (level, goal, player, UI), Layout. Game state: playing / won.
- **internal/player/player.go** — Position, velocity, AABB. Input (A/D/W), gravity, jump (coyote time + jump buffer). Collision is resolved by level; game applies result and calls TryJump.
- **internal/level/level.go** — Level struct (platforms, goal rect, bounds, start, death Y). `FirstLevel(screenW, screenH)` returns level 1. `ResolveCollision(rect, vx, vy)` returns new position, velocity, grounded. `InGoal(rect)` for win check.
- **internal/camera/camera.go** — Camera X,Y in world. `Update(targetX, targetY, levelW, levelH, screenW, screenH)` to follow (lerp) and clamp. `WorldToScreen(wx, wy)` for drawing.

## Conventions

- All game code under `cmd/` and `internal/`. Ebitengine’s Y-down coordinate system; units in pixels. Screen size constants in `internal/game` (640×360).
