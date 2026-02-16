package level

import (
	"image"
)

// PalmTree is a decorative tree position (base X,Y in world coords).
type PalmTree struct {
	X, Y float64
}

// Level holds platform and goal data for one level.
type Level struct {
	Platforms  []image.Rectangle
	Goal       image.Rectangle
	PalmTrees  []PalmTree
	Width      int
	Height     int
	StartX     float64
	StartY     float64
	DeathY     float64 // player dies if Y > DeathY
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

	palms := []PalmTree{
		{X: 30, Y: float64(floorY)},
		{X: 450, Y: float64(floorY)},
		{X: 870, Y: float64(floorY)},
		{X: 1500, Y: float64(floorY)},
		{X: 1900, Y: float64(floorY)},
		{X: 2400, Y: float64(floorY)},
		{X: 4800, Y: float64(floorY)},
	}

	return &Level{
		Platforms: platforms,
		Goal:      goal,
		PalmTrees: palms,
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

// SecondLevel returns the second level: 4 screens wide, different layout, similar difficulty.
func SecondLevel(screenW, screenH int) *Level {
	w := screenW * 4
	h := screenH * 2
	floorY := h - 48

	platforms := []image.Rectangle{
		// No full ground -- floating platforms only, with a floor on the far edges
		// Start platform
		image.Rect(0, floorY, 300, h),

		// Scattered low platforms (screen 1)
		image.Rect(360, floorY-60, 480, floorY),
		image.Rect(540, floorY-130, 660, floorY-60),
		image.Rect(720, floorY-60, 840, floorY),

		// Zigzag section (screen 2)
		image.Rect(920, floorY-100, 1020, floorY-40),
		image.Rect(1060, floorY-200, 1160, floorY-140),
		image.Rect(1200, floorY-100, 1300, floorY-40),
		image.Rect(1340, floorY-200, 1440, floorY-140),
		image.Rect(1480, floorY-300, 1580, floorY-240),

		// Bridge of narrow platforms (screen 3)
		image.Rect(1660, floorY-200, 1720, floorY-160),
		image.Rect(1780, floorY-200, 1840, floorY-160),
		image.Rect(1900, floorY-200, 1960, floorY-160),
		image.Rect(2020, floorY-200, 2080, floorY-160),

		// Descent with wide platforms (screen 3-4)
		image.Rect(2160, floorY-140, 2320, floorY-60),
		image.Rect(2380, floorY-60, 2540, floorY),

		// Climbing section (screen 4)
		image.Rect(2600, floorY-80, 2720, floorY),
		image.Rect(2700, floorY-180, 2820, floorY-100),
		image.Rect(2840, floorY-280, 2960, floorY-200),
		image.Rect(2980, floorY-180, 3100, floorY-100),
		image.Rect(3120, floorY-280, 3240, floorY-200),

		// Ledge walk to the goal
		image.Rect(3320, floorY-220, 3500, floorY-160),

		// Long run to goal
		image.Rect(3560, floorY-120, 3800, floorY),

		// Goal platform
		image.Rect(3860, floorY-200, w, floorY),
	}

	goal := image.Rect(w-120, floorY-220, w-24, floorY-120)
	startX := 80.0
	startY := float64(floorY - 40)

	palms := []PalmTree{
		{X: 100, Y: float64(floorY)},
		{X: 750, Y: float64(floorY)},
		{X: 1250, Y: float64(floorY - 100)},
		{X: 2450, Y: float64(floorY)},
		{X: 3650, Y: float64(floorY)},
		{X: 4000, Y: float64(floorY)},
	}

	return &Level{
		Platforms: platforms,
		Goal:      goal,
		PalmTrees: palms,
		Width:     w,
		Height:    h,
		StartX:    startX,
		StartY:    startY,
		DeathY:    float64(floorY + 100),
	}
}

// ThirdLevel returns the third level: 4 screens wide, harder than level 2.
// Precision platforming with small platforms, wide gaps, and no safety nets.
func ThirdLevel(screenW, screenH int) *Level {
	w := screenW * 4
	h := screenH * 2
	floorY := h - 48

	platforms := []image.Rectangle{
		// Tiny start platform
		image.Rect(0, floorY, 160, h),

		// Screen 1: precision hops on small floating platforms with big gaps
		image.Rect(240, floorY-80, 310, floorY-40),   // 80px up from start
		image.Rect(400, floorY-150, 460, floorY-110),  // 70px up
		image.Rect(540, floorY-80, 600, floorY-40),    // 70px down
		image.Rect(700, floorY-155, 760, floorY-115),  // 75px up (was 100, impossible)
		image.Rect(850, floorY-100, 910, floorY-60),   // 55px down

		// Screen 2: ascending tower -- tiny ledges going up (max 80px per step)
		image.Rect(1000, floorY-60, 1060, floorY-20),   // 40px down
		image.Rect(1100, floorY-140, 1150, floorY-100),  // 80px up
		image.Rect(1200, floorY-220, 1250, floorY-180),  // 80px up
		image.Rect(1290, floorY-295, 1350, floorY-255),  // 75px up (was 90, borderline)
		image.Rect(1400, floorY-370, 1460, floorY-330),  // 75px up (was 90+, impossible)

		// Screen 2-3: high altitude crossing -- very narrow platforms
		image.Rect(1560, floorY-350, 1600, floorY-320),  // 20px down
		image.Rect(1680, floorY-340, 1720, floorY-310),  // 10px down
		image.Rect(1800, floorY-360, 1840, floorY-330),  // 20px up
		image.Rect(1920, floorY-340, 1960, floorY-310),  // 20px down
		image.Rect(2040, floorY-360, 2080, floorY-330),  // 20px up

		// Screen 3: rapid descent -- staircase down with small landings
		image.Rect(2180, floorY-300, 2240, floorY-260),  // 60px down
		image.Rect(2300, floorY-220, 2360, floorY-180),  // 80px down
		image.Rect(2420, floorY-140, 2480, floorY-100),  // 80px down
		image.Rect(2540, floorY-60, 2600, floorY-20),    // 80px down

		// Screen 3-4: the gauntlet -- alternating high/low (max 80px swings)
		image.Rect(2720, floorY-140, 2790, floorY-100),  // 80px up (was 120, impossible)
		image.Rect(2900, floorY-80, 2960, floorY-40),    // 60px down
		image.Rect(3080, floorY-160, 3140, floorY-120),  // 80px up (was 160, impossible)
		image.Rect(3260, floorY-80, 3320, floorY-40),    // 80px down
		image.Rect(3440, floorY-160, 3510, floorY-120),  // 80px up (was 160, impossible)

		// Screen 4: final climb to the goal
		image.Rect(3620, floorY-160, 3680, floorY-120),  // same height
		image.Rect(3740, floorY-240, 3800, floorY-200),  // 80px up (was 100, impossible)
		image.Rect(3880, floorY-310, 3950, floorY-270),  // 70px up

		// Goal platform -- small, must earn it
		image.Rect(4060, floorY-280, 4200, floorY-220),  // 30px down
	}

	goal := image.Rect(4100, floorY-300, 4180, floorY-220)
	startX := 40.0
	startY := float64(floorY - 40)

	palms := []PalmTree{
		{X: 50, Y: float64(floorY)},
		{X: 880, Y: float64(floorY - 60)},
		{X: 1700, Y: float64(floorY - 340)},
		{X: 2500, Y: float64(floorY - 20)},
		{X: 3600, Y: float64(floorY - 120)},
	}

	return &Level{
		Platforms: platforms,
		Goal:      goal,
		PalmTrees: palms,
		Width:     w,
		Height:    h,
		StartX:    startX,
		StartY:    startY,
		DeathY:    float64(floorY + 100),
	}
}

// InGoal returns true if the given rect overlaps the goal area.
func (l *Level) InGoal(rect image.Rectangle) bool {
	return rect.Min.X < l.Goal.Max.X && rect.Max.X > l.Goal.Min.X &&
		rect.Min.Y < l.Goal.Max.Y && rect.Max.Y > l.Goal.Min.Y
}
