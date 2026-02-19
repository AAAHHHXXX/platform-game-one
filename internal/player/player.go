package player

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	Radius        = 14
	Width         = Radius * 2
	Height        = Radius * 2
	MoveSpeed     = 280
	JumpVelocity  = -420
	Gravity       = 980
	CoyoteTimeMax = 0.12
	JumpBufferMax = 0.1
)

// Shape selects the player's visual appearance.
type Shape int

const (
	ShapeCircle   Shape = iota
	ShapeTriangle
	ShapeHexagon
	shapeCount // keep last for cycling
)

// Player represents the controllable character.
type Player struct {
	X, Y       float64
	VX, VY     float64
	Grounded   bool
	Rotation   float64
	Shape      Shape
	CoyoteTime float64
	JumpBuffer float64
}

// New creates a player at the given position.
func New(x, y float64) *Player {
	return &Player{X: x, Y: y, Shape: ShapeCircle}
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

// Update applies input, gravity, and integrates position.
func (p *Player) Update(dt float64) {
	// Toggle shape with Tab (cycle through all shapes)
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		p.Shape = (p.Shape + 1) % shapeCount
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.JumpBuffer = JumpBufferMax
	}
	if p.JumpBuffer > 0 {
		p.JumpBuffer -= dt
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.VX = -MoveSpeed
	} else if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.VX = MoveSpeed
	} else {
		p.VX = 0
	}

	p.VY += Gravity * dt
	p.X += p.VX * dt
	p.Y += p.VY * dt

	p.Rotation += (p.VX * dt) / Radius

	if p.Grounded {
		p.CoyoteTime = CoyoteTimeMax
	} else {
		p.CoyoteTime -= dt
		if p.CoyoteTime < 0 {
			p.CoyoteTime = 0
		}
	}
}

// TryJump applies jump velocity if W was pressed and the player can jump.
func (p *Player) TryJump() {
	if p.JumpBuffer > 0 && (p.Grounded || p.CoyoteTime > 0) {
		p.VY = JumpVelocity
		p.Grounded = false
		p.CoyoteTime = 0
		p.JumpBuffer = 0
	}
}

// Pre-rendered images for each shape.
var (
	circleImg   *ebiten.Image
	triangleImg *ebiten.Image
	hexagonImg  *ebiten.Image
	whitePixel  *ebiten.Image
)

func init() {
	whitePixel = ebiten.NewImage(3, 3)
	whitePixel.Fill(color.White)

	imgSize := Width + 2
	buildCircleImage(imgSize)
	buildTriangleImage(imgSize)
	buildHexagonImage(imgSize)
}

func buildCircleImage(size int) {
	circleImg = ebiten.NewImage(size, size)
	cx := float32(size) / 2
	cy := float32(size) / 2
	r := float32(Radius)

	// Purple body
	vector.DrawFilledCircle(circleImg, cx, cy, r, color.RGBA{R: 0x8a, G: 0x2b, B: 0xe2, A: 0xff}, true)

	drawEyes(circleImg, cx, cy)
}

func buildTriangleImage(size int) {
	triangleImg = ebiten.NewImage(size, size)
	cx := float32(size) / 2
	cy := float32(size) / 2
	r := float32(Radius)

	// Equilateral triangle pointing up, inscribed in the radius
	// Top vertex, bottom-left, bottom-right
	topX := cx
	topY := cy - r
	blX := cx - r*float32(math.Cos(math.Pi/6))
	blY := cy + r*float32(math.Sin(math.Pi/6))
	brX := cx + r*float32(math.Cos(math.Pi/6))
	brY := cy + r*float32(math.Sin(math.Pi/6))

	drawFilledTriangle(triangleImg, topX, topY, blX, blY, brX, brY, color.RGBA{R: 0xd8, G: 0x2a, B: 0x2a, A: 0xff})

	// Eyes positioned near the centroid (slightly above)
	eyeCenterY := cy - r*0.1
	drawEyes(triangleImg, cx, eyeCenterY)
}

func drawEyes(img *ebiten.Image, cx, cy float32) {
	eyeR := float32(4)
	eyeOffX := float32(4.5)
	eyeOffY := float32(-3.5)
	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	vector.DrawFilledCircle(img, cx-eyeOffX, cy+eyeOffY, eyeR, white, true)
	vector.DrawFilledCircle(img, cx+eyeOffX, cy+eyeOffY, eyeR, white, true)

	pupilR := float32(2)
	pupilOffX := float32(5.0)
	pupilOffY := float32(-3.5)
	dark := color.RGBA{R: 0x10, G: 0x10, B: 0x20, A: 0xff}
	vector.DrawFilledCircle(img, cx-pupilOffX, cy+pupilOffY, pupilR, dark, true)
	vector.DrawFilledCircle(img, cx+pupilOffX, cy+pupilOffY, pupilR, dark, true)
}

func buildHexagonImage(size int) {
	hexagonImg = ebiten.NewImage(size, size)
	cx := float32(size) / 2
	cy := float32(size) / 2
	r := float32(Radius)

	// Regular hexagon: 6 vertices evenly spaced, flat-top orientation
	vertices := make([]float32, 12) // 6 pairs of (x, y)
	for i := 0; i < 6; i++ {
		angle := float64(i)*math.Pi/3.0 - math.Pi/6.0 // start at -30 deg for flat top
		vertices[i*2] = cx + r*float32(math.Cos(angle))
		vertices[i*2+1] = cy + r*float32(math.Sin(angle))
	}
	drawFilledPolygon(hexagonImg, vertices, color.RGBA{R: 0xff, G: 0x8c, B: 0x00, A: 0xff})

	drawEyes(hexagonImg, cx, cy)
}

func drawFilledPolygon(dst *ebiten.Image, verts []float32, clr color.Color) {
	var path vector.Path
	path.MoveTo(verts[0], verts[1])
	for i := 2; i < len(verts); i += 2 {
		path.LineTo(verts[i], verts[i+1])
	}
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

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillRuleNonZero,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

func drawFilledTriangle(dst *ebiten.Image, x1, y1, x2, y2, x3, y3 float32, clr color.Color) {
	var path vector.Path
	path.MoveTo(x1, y1)
	path.LineTo(x2, y2)
	path.LineTo(x3, y3)
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

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.FillRuleNonZero,
	}
	dst.DrawTriangles(vs, is, whitePixel, op)
}

// Draw draws the player with the active shape, rolling.
func (p *Player) Draw(screen *ebiten.Image, sx, sy int) {
	var img *ebiten.Image
	switch p.Shape {
	case ShapeTriangle:
		img = triangleImg
	case ShapeHexagon:
		img = hexagonImg
	default:
		img = circleImg
	}

	op := &ebiten.DrawImageOptions{}
	imgSize := float64(img.Bounds().Dx())
	half := imgSize / 2

	op.GeoM.Translate(-half, -half)
	op.GeoM.Rotate(p.Rotation)
	op.GeoM.Translate(half, half)
	op.GeoM.Translate(float64(sx)-1, float64(sy)-1)

	op.Filter = ebiten.FilterLinear
	screen.DrawImage(img, op)
}

// CenterX returns the world X of the player's center.
func (p *Player) CenterX() float64 { return p.X + float64(Width)/2 }

// CenterY returns the world Y of the player's center.
func (p *Player) CenterY() float64 { return p.Y + float64(Height)/2 }
