package game

import (
	"fmt"
	"image/color"

	"platform-game-one/internal/camera"
	"platform-game-one/internal/level"
	"platform-game-one/internal/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	TotalLevels  = 3
)

type gameState int

const (
	statePlaying gameState = iota
	stateWon
)

// Game implements ebiten.Game.
type Game struct {
	player   *player.Player
	level    *level.Level
	camera   *camera.Camera
	state    gameState
	levelNum int // 1-based
}

// New creates a new Game starting at level 1.
func New() *Game {
	lv := level.FirstLevel(ScreenWidth, ScreenHeight)
	pl := player.New(lv.StartX, lv.StartY)
	cam := camera.New()
	return &Game{
		player:   pl,
		level:    lv,
		camera:   cam,
		state:    statePlaying,
		levelNum: 1,
	}
}

func (g *Game) loadLevel(num int) {
	switch num {
	case 1:
		g.level = level.FirstLevel(ScreenWidth, ScreenHeight)
	case 2:
		g.level = level.SecondLevel(ScreenWidth, ScreenHeight)
	case 3:
		g.level = level.ThirdLevel(ScreenWidth, ScreenHeight)
	}
	g.player.Respawn(g.level.StartX, g.level.StartY)
	g.camera = camera.New()
	g.levelNum = num
	g.state = statePlaying
}

// Update runs each tick.
func (g *Game) Update() error {
	if g.state == stateWon {
		return nil
	}
	dt := 1.0 / 60.0

	g.player.Update(dt)
	rect := g.player.Rect()
	nx, ny, nvx, nvy, grounded := g.level.ResolveCollision(rect, g.player.VX, g.player.VY)
	g.player.X, g.player.Y = nx, ny
	g.player.VX, g.player.VY = nvx, nvy
	g.player.Grounded = grounded
	g.player.TryJump()

	// Death: fell below level
	if g.player.Y > g.level.DeathY {
		g.player.Respawn(g.level.StartX, g.level.StartY)
	}
	// Win: reached goal
	if g.level.InGoal(g.player.Rect()) {
		if g.levelNum < TotalLevels {
			g.loadLevel(g.levelNum + 1)
		} else {
			g.state = stateWon
		}
	}

	// Camera follow
	centerX := g.player.CenterX()
	centerY := g.player.CenterY()
	g.camera.Update(centerX, centerY, g.level.Width, g.level.Height, ScreenWidth, ScreenHeight)

	return nil
}

// Draw renders the game.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 0x1a, G: 0x1a, B: 0x2e, A: 0xff})

	for _, plat := range g.level.Platforms {
		sx, sy := g.camera.WorldToScreen(float64(plat.Min.X), float64(plat.Min.Y))
		if sx+plat.Dx() < 0 || sy+plat.Dy() < 0 || sx > ScreenWidth || sy > ScreenHeight {
			continue
		}
		img := ebiten.NewImage(plat.Dx(), plat.Dy())
		img.Fill(color.RGBA{R: 0x4a, G: 0x7c, B: 0x59, A: 0xff})
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(sx), float64(sy))
		screen.DrawImage(img, op)
	}

	// Goal
	goal := g.level.Goal
	sgx, sgy := g.camera.WorldToScreen(float64(goal.Min.X), float64(goal.Min.Y))
	goalImg := ebiten.NewImage(goal.Dx(), goal.Dy())
	goalImg.Fill(color.RGBA{R: 0xea, G: 0xc5, B: 0x4f, A: 0xff})
	opGoal := &ebiten.DrawImageOptions{}
	opGoal.GeoM.Translate(float64(sgx), float64(sgy))
	screen.DrawImage(goalImg, opGoal)

	// Player
	px, py := g.camera.WorldToScreen(g.player.X, g.player.Y)
	g.player.Draw(screen, px, py)

	// HUD
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Level %d / %d   [Tab] switch shape", g.levelNum, TotalLevels))

	if g.state == stateWon {
		ebitenutil.DebugPrint(screen, "\n\n  You beat all three levels! Congratulations!")
	}
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
