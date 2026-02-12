package camera

// Camera holds the top-left position of the visible window in world coordinates.
type Camera struct {
	X, Y float64
}

// New creates a camera at (0,0).
func New() *Camera {
	return &Camera{}
}

// Update moves the camera toward the target (e.g. player center) and clamps to level bounds.
// targetX, targetY is the world position to center on (e.g. player center).
func (c *Camera) Update(targetX, targetY float64, levelW, levelH, screenW, screenH int) {
	// Center the target on screen
	goalX := targetX - float64(screenW)/2
	goalY := targetY - float64(screenH)/2
	// Smooth follow (lerp)
	const speed = 0.12
	c.X += (goalX - c.X) * speed
	c.Y += (goalY - c.Y) * speed
	// Clamp to level
	if c.X < 0 {
		c.X = 0
	}
	if c.Y < 0 {
		c.Y = 0
	}
	maxX := float64(levelW - screenW)
	if maxX > 0 && c.X > maxX {
		c.X = maxX
	}
	maxY := float64(levelH - screenH)
	if maxY > 0 && c.Y > maxY {
		c.Y = maxY
	}
}

// WorldToScreen converts world coordinates to screen coordinates.
func (c *Camera) WorldToScreen(wx, wy float64) (sx, sy int) {
	return int(wx - c.X), int(wy - c.Y)
}
