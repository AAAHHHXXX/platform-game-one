# CONTEXT.md -- Session 1 Handoff

This file captures the full development history and state of the project for the next AI context window.

## Session 1 Summary (Feb 11, 2026)

Built a 2D platformer from scratch in Go with Ebitengine. Started from an empty directory and implemented the full game across multiple feature additions.

## What was built (in order)

### Phase 1: Core game scaffold
- Go module, entry point (`cmd/game/main.go`), game loop (`internal/game/game.go`)
- Player package with W/A/D input, gravity, AABB collision
- Level 1 with platforms defined as `[]image.Rectangle`
- Camera with lerp follow and level-bound clamping
- Win (reach goal zone) and death (fall below level) mechanics

### Phase 2: Bug fixes
- `Fill()` color type: Ebitengine requires `color.RGBA{}`, not hex int literals. Fixed across all files.
- Player auto-jumping: `TryJump()` had a logic bug where `(canJump || jumpBuffer > 0) && canJump` simplified to just `canJump`. Fixed to require `jumpBuffer > 0 && canJump`.

### Phase 3: Rolling circle player with eyes
- Changed player from rectangle to circle (Radius=14, AABB 28x28)
- Added `Rotation` field that accumulates via `VX * dt / Radius` (rolling physics)
- Pre-rendered player image using `vector.DrawFilledCircle` with white eye circles and dark pupils
- Draw applies rotation via GeoM transforms around center

### Phase 4: Player color change
- Changed player body from blue (`#40A0FF`) to purple (`#8A2BE2`)

### Phase 5: Shape toggle, screen resize, level 2
- Added `Shape` type (circle/triangle) with Tab key toggle
- Triangle drawn using `vector.Path` + `DrawTriangles` with a white pixel source
- Doubled screen from 640x360 to 1280x720
- Added `SecondLevel()` with floating platforms, zigzag, bridge, climbing sections
- Game chains levels: beat level 1 -> load level 2 -> beat level 2 -> win

### Phase 6: Orange hexagon + level 3
- Added `ShapeHexagon` with flat-top regular hexagon using `drawFilledPolygon()`
- Tab now cycles through 3 shapes
- Added `ThirdLevel()` with precision platforming: small platforms, ascending tower, high altitude crossing, gauntlet, final climb
- TotalLevels bumped to 3

## Current codebase state

### Git
- **Branch**: `main` (3 commits)
  - **Commits**: `d2d4b0e` Initial -> `09295f3` More levels/characters -> `add commmit SHA` update markdown files
- **Feature branch**: `feature/missionary-megans-updates` (2 commits)
  - **Commits**: `add commit SHA` fix impossible jumps -> `add commit SHA` add sunset/palms/lyrics -> `add commmit SHA` update markdown files
- **Remote**: GitHub (origin)

### File structure
```
platform-game-one/
├── go.mod                      (module platform-game-one, go 1.24, ebiten v2.9.0)
├── go.sum
├── .gitignore
├── README.md
├── CLAUDE.md
├── CONTEXT.md                  (this file)
├── cmd/game/main.go            (17 lines, entry point)
├── internal/
│   ├── game/game.go            (game loop, draw pipeline, level chaining)
│   ├── player/player.go        (3 shapes, rolling, eyes, physics)
│   ├── level/level.go          (3 levels, collision, PalmTree positions)
│   └── camera/camera.go        (44 lines, lerp follow + clamp)
```

### Player physics constants
| Constant | Value | Notes |
|----------|-------|-------|
| Radius | 14 | Circle radius; AABB is 28x28 |
| MoveSpeed | 280 | Pixels/second horizontal |
| JumpVelocity | -420 | Upward (Y-down), applied on jump |
| Gravity | 980 | Pixels/second^2 downward |
| CoyoteTimeMax | 0.12 | Seconds after leaving ground where jump still works |
| JumpBufferMax | 0.1 | Seconds before landing where W press is remembered |
| **Max jump height** | **~90px** | Derived: v^2/(2g) = 420^2/(2*980) |
| **Max horizontal (same height)** | **~240px** | Derived: MoveSpeed * 2 * (JumpVelocity/Gravity) |

### Level design constraints
- All upward jumps in level data MUST have height difference <= 80px (with 90px max, this leaves margin)
- Level width = `screenW * 4` (5120px at 1280 screen width)
- Level height = `screenH * 2` (1440px)
- floorY = levelHeight - 48 = 1392

## Known issues / areas for improvement

1. **Platform images created every frame in Draw()**: Each `ebiten.NewImage()` call in the platform drawing loop allocates a new texture per frame. Should be cached (pre-render one per unique size, or use a single white image + scale via GeoM).
2. **No sound effects or music**.
3. **No animations** (e.g., squash/stretch on land, particle effects).
4. **Win screen is minimal** -- just debug text. Could show a proper victory screen with restart option.
5. **No restart key** -- once you win, the game freezes. Could add R to restart.
6. **CLAUDE.md and README.md may be out of date** -- they still reference single-level, 640x360, etc. Update them when making changes.
7. **Level 3 may still have tight/borderline jumps** -- playtest and adjust if needed.
8. **No title screen or level select**.

## Potential future features (discussed but not implemented)
- Moving platforms, one-way platforms, breakable blocks
- Sprites/pixel art instead of vector shapes
- Loading levels from JSON/Tiled files
- Sound effects and music
- Multiple lives / score tracking
- Enemies (was explicitly out of scope for v1)
