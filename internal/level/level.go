package level

import (
	"image"
)

// Level holds platform and goal data for one level.
type Level struct {
	Platforms   []image.Rectangle
	Goal        image.Rectangle
	Width       int
	Height      int
	StartX      float64
	StartY      float64
	DeathY      float64 // player dies if Y > DeathY
}

// FirstLevel returns the first level: 4 screens wide, complex layout.
func FirstLevel(screenW, screenH int) *Level {
	w := screenW * 4
	h := screenH * 2
	floorY := h - 48

	// Platforms: (minX, minY, maxX, maxY) using image.Rect
	platforms := []image.Rectangle{
		// Ground
		image.Rect(0, floorY, w, h),
		// First screen: starting area and first obstacles
		image.Rect(80, floorY - 40, 200, floorY),
		image.Rect(240, floorY - 120, 360, floorY - 40),
		image.Rect(400, floorY - 80, 520, floorY),
		image.Rect(560, floorY - 160, 640, floorY - 80),
		// Gap then platform
		image.Rect(680, floorY - 200, 800, floorY - 80),
		image.Rect(820, floorY - 120, 960, floorY),
		// Stairs up
		image.Rect(1000, floorY - 40, 1100, floorY),
		image.Rect(1080, floorY - 100, 1180, floorY - 40),
		image.Rect(1160, floorY - 160, 1260, floorY - 100),
		image.Rect(1240, floorY - 220, 1340, floorY - 160),
		// Narrow ledge
		image.Rect(1380, floorY - 260, 1480, floorY - 200),
		// Drop and climb
		image.Rect(1540, floorY - 180, 1660, floorY - 100),
		image.Rect(1620, floorY - 280, 1760, floorY - 180),
		// Final approach
		image.Rect(1820, floorY - 120, 1980, floorY),
		image.Rect(2000, floorY - 200, 2120, floorY - 80),
		image.Rect(2180, floorY - 280, 2300, floorY - 160),
		// Goal platform
		image.Rect(2360, floorY - 200, w, floorY),
	}

	goal := image.Rect(w-120, floorY-220, w-24, floorY-120)
	startX := 64.0
	startY := float64(floorY - 40 - 32) // above first small platform

	return &Level{
		Platforms: platforms,
		Goal:      goal,
		Width:     w,
		Height:    h,
		StartX:    startX,
		StartY:    startY,
		DeathY:    float64(floorY + 100),
	}
}

// ResolveCollision takes the player's current rect and velocity, resolves collisions
// with all platforms, and returns the new position (as min X,Y of rect), new velocity,
// and whether the player is grounded. Multiple passes ensure we don't stay stuck.
func (l *Level) ResolveCollision(rect image.Rectangle, vx, vy float64) (newX, newY float64, newVX, newVY float64, grounded bool) {
	newX = float64(rect.Min.X)
	newY = float64(rect.Min.Y)
	newVX = vx
	newVY = vy
	const maxPasses = 4
	for pass := 0; pass < maxPasses; pass++ {
		anyResolved := false
		for _, plat := range l.Platforms {
			rect = image.Rect(int(newX), int(newY), int(newX)+rect.Dx(), int(newY)+rect.Dy())
			if rect.Min.X >= plat.Max.X || rect.Max.X <= plat.Min.X ||
				rect.Min.Y >= plat.Max.Y || rect.Max.Y <= plat.Min.Y {
				continue
			}
			overlapLeft := float64(rect.Max.X - plat.Min.X)
			overlapRight := float64(plat.Max.X - rect.Min.X)
			overlapTop := float64(rect.Max.Y - plat.Min.Y)
			overlapBottom := float64(plat.Max.Y - rect.Min.Y)
			minOverlap := overlapLeft
			axis := 0
			if overlapRight < minOverlap {
				minOverlap = overlapRight
				axis = 0
			}
			if overlapTop < minOverlap {
				minOverlap = overlapTop
				axis = 1
			}
			if overlapBottom < minOverlap {
				minOverlap = overlapBottom
				axis = 1
			}
			if axis == 0 {
				if overlapLeft < overlapRight {
					newX = float64(plat.Min.X) - float64(rect.Dx())
					newVX = 0
				} else {
					newX = float64(plat.Max.X)
					newVX = 0
				}
				anyResolved = true
			} else {
				if overlapTop < overlapBottom {
					newY = float64(plat.Min.Y) - float64(rect.Dy())
					newVY = 0
					grounded = true
				} else {
					newY = float64(plat.Max.Y)
					newVY = 0
				}
				anyResolved = true
			}
		}
		if !anyResolved {
			break
		}
	}
	return newX, newY, newVX, newVY, grounded
}

// InGoal returns true if the given rect overlaps the goal area.
func (l *Level) InGoal(rect image.Rectangle) bool {
	return rect.Min.X < l.Goal.Max.X && rect.Max.X > l.Goal.Min.X &&
		rect.Min.Y < l.Goal.Max.Y && rect.Max.Y > l.Goal.Min.Y
}
