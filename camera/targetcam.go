package camera

import (
	"log"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/core/entity"
)

var _ CameraI = (*TargetCamera)(nil)

// TODO: TargetCamera does a lot of math for its basic calculations (like the Up function). Consider making more members
// private so camera rotation can be cached.

// TargetCamera is a camera that is always positioned at an offset from the target entity.
// Target Position + Offset = Camera Position.
// Zoomer can modify the length of the offset.
// The camera always looks toward the target with the provided Up vector determining what orientation the viewport has.
// Create one using NewTargetCamera unless you know what you're doing.
type TargetCamera struct {
	Camera

	// Target is an entity which the TargetCamera follows. The camera always faces it and stays TargetOffset away from it.
	Target entity.Target
	// TargetOffset determines where the camera is positioned in relation to the target.
	// Camera.Target.Position + Camera.TargetOffset == Camera.Position
	TargetOffset mgl32.Vec3
	// Zoomer handles camera zoom.
	Zoomer zoom.Zoom
}

// ModelView returns a matrix used to transform from model space to camera coordinates.
// http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
func (c *TargetCamera) ModelView() mgl32.Mat4 {
	// Note that ModelView must be an override in order to call TargetCamera's Forward and Up methods.
	// If this method didn't exist, a call to targetCam.ModelView() would result in the TargetCamera's embedded
	// camera ModelView() being called, which would in turn call the TargetCamera.Camera.Entity's Forward and Up
	// methods.

	if (c.Up() == mgl32.Vec3{0, 0, 0}) {
		log.Println("invalid ModelView: up vector is (0,0,0)")
	}
	if (c.Forward() == mgl32.Vec3{0, 0, 0}) {
		log.Println("invalid ModelView: forward vector is (0,0,0)")
	}
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Forward()), c.Up())
}

// ProjectionOrthographic returns a matrix used to transform from camera space to screen space.
func (c *TargetCamera) ProjectionOrthographic(width, height float32) mgl32.Mat4 {
	// Since distance from target doesn't do a "zoom" effect in an orthographic projection, simulate one
	// by changing how wide the view is.
	zoomPercent := c.GetCurrentZoomPercent()
	return c.Camera.ProjectionOrthographic(width/zoomPercent, height/zoomPercent)
}

func (c *TargetCamera) Update(delta time.Duration) {
	c.Camera.Update(delta)
	if c.Zoomer != nil {
		c.Zoomer.Update()
	}
	// Adjust the distance from camera to target by the amount of zoom.
	// A zoom of 3 means everything should be 3 times as large, so the distance from target to camera should be 1/3 the default.
	offset := c.TargetOffset.Mul(1.0 / c.GetCurrentZoomPercent())
	c.Position = c.Target.GetPosition().Add(offset)
	c.Rotation = mgl32.QuatLookAtV(c.Position, c.Target.GetPosition(), c.Up())
}

func NewTargetCamera(target entity.Target, offset mgl32.Vec3) *TargetCamera {
	c := &TargetCamera{
		Camera:       *NewCamera(),
		Target:       target,
		TargetOffset: offset,
	}
	// The camera should always face toward the target.
	c.Rotation = mgl32.QuatLookAtV(c.Position, c.Target.GetPosition(), c.Up())
	return c
}

// Up returns a vector perpendicular to both the Forward and Right vectors.
// Given the perpendicular constraint, it's the most similar vector to
// entity.Up as possible. In the event that the the camera's Forward
// vector is the same as entity.Up, there is no single "most similar to Up"
// so it defaults to entity.Right.
func (c *TargetCamera) Up() mgl32.Vec3 {
	if c.Forward() == entity.Up {
		return entity.Right
	}
	return c.Right().Cross(c.Forward()).Normalize()
}

func (c *TargetCamera) Forward() mgl32.Vec3 {
	return c.TargetOffset.Mul(-1).Normalize()
}

func (c *TargetCamera) Right() mgl32.Vec3 {
	if c.Forward() == entity.Right {
		return entity.Forward
	}
	return c.Forward().Cross(entity.Up)
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
