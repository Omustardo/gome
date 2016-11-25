package camera

import "github.com/go-gl/mathgl/mgl32"

type Camera interface {
	ModelView() mgl32.Mat4
	Projection(width, height float32) mgl32.Mat4
	Near() float32
	Far() float32
	Update()
	Position() mgl32.Vec3

	// GetCurrentZoomPercent returns how much the camera is zoomed in or out. If the camera doesn't implement zoom
	// functions then this always returns 1.
	// Values smaller than 1 indicate zoomed out. e.g. 0.25 means all objects are 25% of their original size.
	// Values larger than 1 indicate zoomed in. e.g. 3.0 means all objects appear 3 times larger than their original size.
	GetCurrentZoomPercent() float32

	// ScreenToWorldCoord2D returns the world coordinates of a point on the screen.
	ScreenToWorldCoord2D(screenPoint mgl32.Vec2, windowWidth, windowHeight int) mgl32.Vec2
}
