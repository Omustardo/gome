package camera

import "github.com/go-gl/mathgl/mgl32"

// Compile time check that FreeCamera implements Camera
var _ Camera = (*FreeCamera)(nil)

// FreeCamera is an orthographic camera that is not attached to any player. It can be moved by modifying the Pos.
type FreeCamera struct {
	Pos mgl32.Vec3
}

func (c *FreeCamera) ModelView() mgl32.Mat4 {
	return mgl32.LookAt(
		c.Pos[0], c.Pos[1], c.Pos[2], // Camera Position
		c.Pos[0], c.Pos[1], -1, // Target Position. Looking down on Z.
		0, 1, 0) // Up vector
}

func (c *FreeCamera) Projection(width, height float32) mgl32.Mat4 {
	return mgl32.Ortho(-width/2, width/2,
		-height/2, height/2,
		c.Near(), c.Far())
}

func (c *FreeCamera) Update() {
	// Enforce bounds. This camera is always looking down onto the XY plane.
	if c.Pos[2] <= 0 {
		c.Pos[2] = 0
	}
}

func (c *FreeCamera) Near() float32 {
	return 0.1
}

func (c *FreeCamera) Far() float32 {
	return 100
}

func (c *FreeCamera) Position() mgl32.Vec3 {
	return c.Pos
}

func (c *FreeCamera) GetCurrentZoomPercent() float32 {
	return 1
}

func (c *FreeCamera) ScreenToWorldCoord2D(screenPoint mgl32.Vec2, windowSize [2]int) mgl32.Vec2 {
	return mgl32.Vec2{
		c.Pos.X() + (screenPoint.X() - float32(windowSize[0])/2),
		c.Pos.Y() - (screenPoint.Y() - float32(windowSize[1])/2),
	}
}
