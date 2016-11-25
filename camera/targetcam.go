package camera

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/entity"
)

var _ Camera = (*TargetCamera)(nil)

// TargetCamera is an orthographic camera that is always locked to an entity.
type TargetCamera struct {
	Target entity.Entity
	Pos    mgl32.Vec3 // This should always reflect the target's position, but always with a Z value > 0
	zoomer zoom.Zoom
}

func NewTargetCamera(target entity.Entity, zoomer zoom.Zoom) Camera {
	p := &TargetCamera{
		Target: target,
		zoomer: zoomer,
	}
	p.Update()
	p.Pos[2] = 1 // Not necessary with the bounds enforcing in Update(), but nice to be safe.
	return p
}

func (c *TargetCamera) ModelView() mgl32.Mat4 {
	targetPos := c.Target.Center()
	return mgl32.LookAt(
		targetPos.X(), targetPos.Y(), 1, // Camera Position. Always above target.
		targetPos.X(), targetPos.Y(), 0, // Target Position.
		0, 1, 0) // Up vector // TODO: For a camera that has Screen Up always as the direction the entity is facing, I think we just need to modify this line.
}

func (c *TargetCamera) Projection(width, height float32) mgl32.Mat4 {
	zoomPercent := c.GetCurrentZoomPercent()
	width /= zoomPercent
	height /= zoomPercent
	return mgl32.Ortho(-width/2, width/2,
		-height/2, height/2,
		c.Near(), c.Far())
}

func (c *TargetCamera) Update() {
	c.Pos = c.Target.Center()
	if c.zoomer != nil {
		c.zoomer.Update()
	}
}

func (c *TargetCamera) Near() float32 {
	return 0.1
}

func (c *TargetCamera) Far() float32 {
	return 1000
}

func (c *TargetCamera) Position() mgl32.Vec3 {
	pos := c.Target.Center()
	pos[2] = 1
	return pos
}

// ScreenToWorldCoord2D returns the world coordinates of a point on the screen.
// This depends on the camera always looking directly down onto the XY-plane. e.g. camera position has positive Z.
// The screen space coordinate is expected in the coordinate system where the top left corner is (0,0), Y increases down, and X increases right.
func (c *TargetCamera) ScreenToWorldCoord2D(screenPoint mgl32.Vec2, windowWidth, windowHeight int) mgl32.Vec2 {
	zoomPercent := c.GetCurrentZoomPercent()
	return mgl32.Vec2{
		c.Pos.X() + (screenPoint.X()-float32(windowWidth)/2)/zoomPercent,
		c.Pos.Y() - (screenPoint.Y()-float32(windowHeight)/2)/zoomPercent,
	}
}

func (c *TargetCamera) GetCurrentZoomPercent() float32 {
	if c.zoomer == nil {
		return 1
	}
	zoomPercent := c.zoomer.GetCurrentPercent()
	if zoomPercent <= 0 {
		log.Printf("Invalid camera zoom: %v. Using default", zoomPercent)
		zoomPercent = 1.0
	}
	return zoomPercent
}
