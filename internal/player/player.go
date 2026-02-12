package player

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Width          = 24
	Height         = 32
	MoveSpeed      = 280
	JumpVelocity   = -420
	Gravity        = 980
	CoyoteTimeMax  = 0.12
	JumpBufferMax  = 0.1
)

// Player represents the controllable character.
type Player struct {
	X, Y     float64
	VX, VY   float64
	Grounded bool
	// Coyote time: allow jump shortly after leaving ground
	CoyoteTime float64
	// Jump buffer: allow jump shortly before landing
	JumpBuffer float64
}

// New creates a player at the given position.
func New(x, y float64) *Player {
	return &Player{X: x, Y: y}
}

// Rect returns the axis-aligned bounding box in world coordinates.
func (p *Player) Rect() image.Rectangle {
	return image.Rect(
		int(p.X), int(p.Y),
		int(p.X)+Width, int(p.Y)+Height,
	)
}

// Respawn moves the player to the start position and zeroes velocity.
func (p *Player) Respawn(x, y float64) {
	p.X, p.Y = x, y
	p.VX, p.VY = 0, 0
	p.Grounded = false
	p.CoyoteTime = 0
	p.JumpBuffer = 0
}

// Update applies input, gravity, and returns the new AABB and velocity for collision.
// dt is delta time in seconds. After calling Update, the game should call level.ResolveCollision
// and then apply the returned position/velocity back to the player.
func (p *Player) Update(dt float64) {
	// Jump buffer: if W was pressed recently, jump when we land
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		p.JumpBuffer = JumpBufferMax
	}
	if p.JumpBuffer > 0 {
		p.JumpBuffer -= dt
	}

	// Horizontal input
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.VX = -MoveSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.VX = MoveSpeed
	} else {
		p.VX = 0
	}

	// Gravity
	p.VY += Gravity * dt

	// Integrate position (will be corrected by collision)
	p.X += p.VX * dt
	p.Y += p.VY * dt

	// Coyote time: count down when in air
	if p.Grounded {
		p.CoyoteTime = CoyoteTimeMax
	} else {
		p.CoyoteTime -= dt
		if p.CoyoteTime < 0 {
			p.CoyoteTime = 0
		}
	}
}

// TryJump applies jump velocity if conditions are met. Call after collision resolution
// so Grounded is set. Also consumes jump buffer.
func (p *Player) TryJump() {
	canJump := p.Grounded || p.CoyoteTime > 0
	if (canJump || p.JumpBuffer > 0) && (p.Grounded || p.CoyoteTime > 0) {
		p.VY = JumpVelocity
		p.Grounded = false
		p.CoyoteTime = 0
		p.JumpBuffer = 0
	}
}

var playerImg *ebiten.Image

func init() {
	playerImg = ebiten.NewImage(Width, Height)
	playerImg.Fill(color.RGBA{R: 0x40, G: 0xa0, B: 0xff, A: 0xff})
}

// Draw draws the player onto the given image at screen coordinates (sx, sy).
func (p *Player) Draw(screen *ebiten.Image, sx, sy int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sx), float64(sy))
	screen.DrawImage(playerImg, op)
}
