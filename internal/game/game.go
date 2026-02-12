package game

import (
	"image/color"

	"platform-game-one/internal/camera"
	"platform-game-one/internal/level"
	"platform-game-one/internal/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 360
)

type gameState int

const (
	statePlaying gameState = iota
	stateWon
)

// Game implements ebiten.Game.
type Game struct {
	player *player.Player
	level  *level.Level
	camera *camera.Camera
	state  gameState
}

// New creates a new Game.
func New() *Game {
	lv := level.FirstLevel(ScreenWidth, ScreenHeight)
	pl := player.New(lv.StartX, lv.StartY)
	cam := camera.New()
	return &Game{
		player: pl,
		level:  lv,
		camera: cam,
		state:  statePlaying,
	}
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
		g.state = stateWon
	}

	// Camera follow
	centerX := g.player.X + player.Width/2
	centerY := g.player.Y + player.Height/2
	g.camera.Update(centerX, centerY, g.level.Width, g.level.Height, ScreenWidth, ScreenHeight)

	return nil
}

// Draw renders the game.
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear
	screen.Fill(color.RGBA{R: 0x1a, G: 0x1a, B: 0x2e, A: 0xff})

	// Draw level (platforms and goal) in screen space
	for _, plat := range g.level.Platforms {
		sx, sy := g.camera.WorldToScreen(float64(plat.Min.X), float64(plat.Min.Y))
		// Cull off-screen
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

	if g.state == stateWon {
		ebitenutil.DebugPrint(screen, "You win! Reach the end.")
	}
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
