package main

import (
	"platform-game-one/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.New()
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	ebiten.SetWindowTitle("Platformer - Level One")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
