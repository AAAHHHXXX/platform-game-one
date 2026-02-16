package game

import (
	"fmt"
	"image/color"
	"math"

	"platform-game-one/internal/camera"
	"platform-game-one/internal/level"
	"platform-game-one/internal/player"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

const pinaColadaLyrics = `I was tired of my lady, we'd been together too long ` +
	`Like a worn out recording of a favorite song ` +
	`So while she lay there sleeping, I read the paper in bed ` +
	`And in the personal columns, there was this letter I read ` +
	`  ~  ` +
	`"If you like pina coladas, and getting caught in the rain ` +
	`If you're not into yoga, if you have half a brain ` +
	`If you like making love at midnight, in the dunes of the cape ` +
	`I'm the love that you've looked for, write to me and escape" ` +
	`  ~  ` +
	`I didn't think about my lady, I know that sounds kind of mean ` +
	`But me and my old lady, had fallen into the same old dull routine ` +
	`So I wrote to the paper, took out a personal ad ` +
	`And though I'm nobody's poet, I thought it wasn't half bad ` +
	`  ~  ` +
	`"Yes, I like pina coladas, and getting caught in the rain ` +
	`I'm not much into health food, I am into champagne ` +
	`I've got to meet you by tomorrow noon, and cut through all this red tape ` +
	`At a bar called O'Malley's, where we'll plan our escape" ` +
	`  ~  ` +
	`So I waited with high hopes, then she walked in the place ` +
	`I knew her smile in an instant, I knew the curve of her face ` +
	`It was my own lovely lady, and she said, "Oh, it's you" ` +
	`And we laughed for a moment, and I said, "I never knew" ` +
	`  ~  ` +
	`"That you liked pina coladas, and getting caught in the rain ` +
	`And the feel of the ocean, and the taste of champagne ` +
	`If you like making love at midnight, in the dunes of the cape ` +
	`You're the love that I've looked for, come with me and escape"` +
	`          `

// Pre-rendered assets
var (
	bgImage    *ebiten.Image
	palmImage  *ebiten.Image
	whitePixel *ebiten.Image
)

func init() {
	whitePixel = ebiten.NewImage(3, 3)
	whitePixel.Fill(color.White)
	buildBackground()
	buildPalmTree()
}

func buildBackground() {
	bgImage = ebiten.NewImage(ScreenWidth, ScreenHeight)

	// Sunset sky gradient
	horizonY := ScreenHeight * 55 / 100
	for y := 0; y < ScreenHeight; y++ {
		var r, g, b uint8
		if y < horizonY {
			// Sky: dark purple at top → warm orange at horizon
			t := float64(y) / float64(horizonY)
			r = uint8(25 + t*210)
			g = uint8(10 + t*100)
			b = uint8(80 - t*30)
		} else if y < horizonY+40 {
			// Horizon glow band
			t := float64(y-horizonY) / 40.0
			r = uint8(255 - t*50)
			g = uint8(140 - t*80)
			b = uint8(50 + t*30)
		} else if y < ScreenHeight*80/100 {
			// Water
			t := float64(y-horizonY-40) / float64(ScreenHeight*80/100-horizonY-40)
			r = uint8(20 + t*10)
			g = uint8(50 + t*20)
			b = uint8(100 - t*30)
		} else {
			// Sand
			t := float64(y-ScreenHeight*80/100) / float64(ScreenHeight*20/100)
			r = uint8(194 - t*30)
			g = uint8(170 - t*20)
			b = uint8(120 - t*20)
		}
		for x := 0; x < ScreenWidth; x++ {
			bgImage.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 0xff})
		}
	}

	// Sun disc near horizon
	sunCX, sunCY, sunR := float32(ScreenWidth*3/4), float32(horizonY-10), float32(50)
	vector.DrawFilledCircle(bgImage, sunCX, sunCY, sunR, color.RGBA{R: 0xff, G: 0xd7, B: 0x40, A: 0xcc}, true)
	vector.DrawFilledCircle(bgImage, sunCX, sunCY, sunR-8, color.RGBA{R: 0xff, G: 0xe8, B: 0x80, A: 0xdd}, true)
}

func buildPalmTree() {
	w, h := 80, 140
	palmImage = ebiten.NewImage(w, h)
	cx := float32(w) / 2

	// Trunk: slightly curved brown rectangle
	trunkW := float32(8)
	trunkBot := float32(h)
	trunkTop := float32(h - 100)
	vector.DrawFilledRect(palmImage, cx-trunkW/2, trunkTop, trunkW, trunkBot-trunkTop,
		color.RGBA{R: 0x8b, G: 0x5e, B: 0x3c, A: 0xff}, true)

	// Fronds: 5 arcs radiating from the trunk top
	frondColor := color.RGBA{R: 0x22, G: 0x8b, B: 0x22, A: 0xff}
	frondLen := float32(45)
	angles := []float64{-80, -40, 0, 40, 80}
	for _, deg := range angles {
		rad := deg * math.Pi / 180
		endX := cx + frondLen*float32(math.Sin(rad))
		endY := trunkTop - frondLen*float32(math.Cos(rad))*0.6
		midX := cx + frondLen*0.5*float32(math.Sin(rad))
		midY := trunkTop - frondLen*0.5*float32(math.Cos(rad))*0.6 - 8

		// Draw thick frond as two triangles forming a leaf shape
		drawFrond(palmImage, cx, trunkTop, midX, midY, endX, endY, frondColor)
	}
	// Coconuts
	vector.DrawFilledCircle(palmImage, cx-5, trunkTop+4, 4, color.RGBA{R: 0x6b, G: 0x4e, B: 0x2c, A: 0xff}, true)
	vector.DrawFilledCircle(palmImage, cx+5, trunkTop+6, 4, color.RGBA{R: 0x6b, G: 0x4e, B: 0x2c, A: 0xff}, true)
}

func drawFrond(dst *ebiten.Image, x1, y1, midX, midY, x2, y2 float32, clr color.Color) {
	// Leaf shape: two triangles sharing the centerline
	thickness := float32(6)
	// Perpendicular offset
	dx := x2 - x1
	dy := y2 - y1
	ln := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if ln < 1 {
		return
	}
	px := -dy / ln * thickness / 2
	py := dx / ln * thickness / 2

	var path vector.Path
	path.MoveTo(x1, y1)
	path.LineTo(midX+px, midY+py)
	path.LineTo(x2, y2)
	path.LineTo(midX-px, midY-py)
	path.Close()

	cr, cg, cb, ca := clr.RGBA()
	rf := float32(cr) / 0xffff
	gf := float32(cg) / 0xffff
	bf := float32(cb) / 0xffff
	af := float32(ca) / 0xffff

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = rf
		vs[i].ColorG = gf
		vs[i].ColorB = bf
		vs[i].ColorA = af
	}
	op := &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillRuleNonZero}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

// Game implements ebiten.Game.
type Game struct {
	player      *player.Player
	level       *level.Level
	camera      *camera.Camera
	state       gameState
	levelNum    int
	lyricsScroll float64
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

	// Scroll lyrics marquee
	g.lyricsScroll += 60.0 * dt // pixels per second

	g.player.Update(dt)
	rect := g.player.Rect()
	nx, ny, nvx, nvy, grounded := g.level.ResolveCollision(rect, g.player.VX, g.player.VY)
	g.player.X, g.player.Y = nx, ny
	g.player.VX, g.player.VY = nvx, nvy
	g.player.Grounded = grounded
	g.player.TryJump()

	if g.player.Y > g.level.DeathY {
		g.player.Respawn(g.level.StartX, g.level.StartY)
	}
	if g.level.InGoal(g.player.Rect()) {
		if g.levelNum < TotalLevels {
			g.loadLevel(g.levelNum + 1)
		} else {
			g.state = stateWon
		}
	}

	centerX := g.player.CenterX()
	centerY := g.player.CenterY()
	g.camera.Update(centerX, centerY, g.level.Width, g.level.Height, ScreenWidth, ScreenHeight)

	return nil
}

// Draw renders the game.
func (g *Game) Draw(screen *ebiten.Image) {
	// Sunset background with parallax
	bgOp := &ebiten.DrawImageOptions{}
	parallaxX := -g.camera.X * 0.15
	bgOp.GeoM.Translate(parallaxX, 0)
	screen.DrawImage(bgImage, bgOp)
	// Draw a second copy for seamless wrap
	bgOp2 := &ebiten.DrawImageOptions{}
	bgOp2.GeoM.Translate(parallaxX+float64(ScreenWidth), 0)
	screen.DrawImage(bgImage, bgOp2)

	// Palm trees (behind platforms, in front of background)
	for _, pt := range g.level.PalmTrees {
		pw := float64(palmImage.Bounds().Dx())
		ph := float64(palmImage.Bounds().Dy())
		sx, sy := g.camera.WorldToScreen(pt.X-pw/2, pt.Y-ph)
		if sx+int(pw) < 0 || sx > ScreenWidth {
			continue
		}
		pop := &ebiten.DrawImageOptions{}
		pop.GeoM.Translate(float64(sx), float64(sy))
		screen.DrawImage(palmImage, pop)
	}

	// Platforms (magenta pink)
	for _, plat := range g.level.Platforms {
		sx, sy := g.camera.WorldToScreen(float64(plat.Min.X), float64(plat.Min.Y))
		if sx+plat.Dx() < 0 || sy+plat.Dy() < 0 || sx > ScreenWidth || sy > ScreenHeight {
			continue
		}
		img := ebiten.NewImage(plat.Dx(), plat.Dy())
		img.Fill(color.RGBA{R: 0xe0, G: 0x3e, B: 0x8c, A: 0xff})
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

	// Lyrics marquee at the top
	textWidth := float64(len(pinaColadaLyrics)) * 6 // approximate char width
	scrollX := float64(ScreenWidth) - math.Mod(g.lyricsScroll, textWidth+float64(ScreenWidth))
	ebitenutil.DebugPrintAt(screen, pinaColadaLyrics, int(scrollX), 2)

	// HUD below lyrics
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Level %d / %d   [Tab] switch shape", g.levelNum, TotalLevels), 4, 18)

	if g.state == stateWon {
		ebitenutil.DebugPrintAt(screen, "You beat all three levels! Congratulations!", ScreenWidth/2-150, ScreenHeight/2)
	}
}

// Layout returns the logical screen size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}
