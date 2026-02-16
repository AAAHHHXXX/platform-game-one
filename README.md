# Platformer

A single-player 2D platformer written in Go using [Ebitengine](https://ebitengine.org/). Navigate 3 increasingly challenging levels to win.

## Controls

- **W** -- Jump
- **A** -- Move left
- **D** -- Move right
- **Tab** -- Cycle player shape (purple circle / red triangle / orange hexagon)

## Build and run

```bash
go mod tidy
go run ./cmd/game
```

Build a binary:

```bash
go build -o platformer ./cmd/game
./platformer
```

## Features

- 3 levels, each 4 screens wide, with progressively harder platforming
- 3 playable shapes (all with eyes, all roll when moving)
- Beach sunset background with parallax scrolling
- Decorative palm trees
- Scrolling song lyrics at the top of the screen
- Camera follows the player and stays within level bounds
- Fall off the bottom to respawn at the start; reach the gold goal to advance
