# Platformer - Level One

A single-player 2D platformer written in Go using [Ebitengine](https://ebitengine.org/). Reach the end of the level to win.

## Controls

- **W** — Jump  
- **A** — Move left  
- **D** — Move right  

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

- One level, four screens wide, with a challenging layout (platforms, gaps, stairs).
- Camera follows the player and stays within level bounds.
- Fall off the bottom to respawn at the start; touch the goal (gold) to win.
