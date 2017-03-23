package camera

import (
	"log"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/core/entity"
)

var _ CameraI = (*TargetCamera)(nil)

// TargetCamera is a camera that is always positioned at an offset from the target entity. Zoomer can modify the length
// of the offset. The camera always looks toward the target with the provided Up vector determining what orientation
// the viewport has.
type TargetCamera struct {
	Camera

	Target entity.Target
	// TargetOffset determines where the camera is positioned in relation to the target.
	TargetOffset mgl32.Vec3
	// Zoomer handles camera zoom.
	Zoomer zoom.Zoom
}

func (c *TargetCamera) ModelView() mgl32.Mat4 {
	if (c.Up() == mgl32.Vec3{0, 0, 0}) {
		log.Println("invalid ModelView: up vector is (0,0,0)")
	}
	if (c.Forward() == mgl32.Vec3{0, 0, 0}) {
		log.Println("invalid ModelView: forward vector is (0,0,0)")
	}
	return mgl32.LookAtV(c.Entity.Position, c.Target.GetPosition(), c.Up())
}

func (c *TargetCamera) Update(delta time.Duration) {
	c.Camera.Update(delta)
	if c.Zoomer != nil {
		c.Zoomer.Update()
	}

	// Adjust the distance from camera to target by the amount of zoom.
	// A zoom of 3 means everything should be 3 times as large, so the distance from target to camera should be 1/3 the default.
	offset := c.TargetOffset.Mul(1.0 / c.GetCurrentZoomPercent())
	c.Entity.Position = c.Target.GetPosition().Add(offset)
	// TODO: I'm unsure if this rotation is being set properly. Easiest test would be to render the camera entity in a model
	// so it's obvious what the rotation is.
	c.Entity.Rotation = mgl32.QuatLookAtV(c.Position, c.Position.Sub(c.TargetOffset), c.Up()) // mgl32.QuatRotate(0, c.TargetOffset) // mgl32.QuatLookAtV(c.Entity.Position, c.Entity.Position.Add(c.TargetOffset), c.Up())
}

func NewTargetCamera(target entity.Target, offset mgl32.Vec3) *TargetCamera {
	return &TargetCamera{
		Camera:       *NewCamera(),
		Target:       target,
		TargetOffset: offset,
	}
}

func (c *TargetCamera) Forward() mgl32.Vec3 {
	return c.TargetOffset
}
func (c *TargetCamera) Up() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0} // TODO: How do I rotate this to work with the Forward vector?
}

func (c *TargetCamera) GetCurrentZoomPercent() float32 {
	if c.Zoomer == nil {
		return 1
	}
	zoomPercent := c.Zoomer.GetCurrentPercent()
	if zoomPercent <= 0 {
		log.Printf("Invalid camera zoom: %v. Using default", zoomPercent)
		zoomPercent = 1.0
	}
	return zoomPercent
}
