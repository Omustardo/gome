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

	// Target is an entity which the TargetCamera follows. The camera always faces it and stays TargetOffset away from it.
	Target entity.Target
	// TargetOffset determines where the camera is positioned in relation to the target.
	// Camera.Target.Position + Camera.TargetOffset == Camera.Position
	TargetOffset mgl32.Vec3
	// Zoomer handles camera zoom.
	Zoomer zoom.Zoom

	// prevOffset keeps track of the previous TargetOffset. This allows us to avoid recomputing some quaternions on every call to Update.
	prevOffset mgl32.Vec3
}

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
	// Only modify the camera rotation if the offset has changed.
	// TODO: I added this because if I update Rotation every Update it causes the screen to flicker. I'm not sure why. This is a temporary workaround.
	if c.prevOffset != c.TargetOffset {
		c.prevOffset = c.TargetOffset
		c.Rotation = mgl32.QuatLookAtV(c.Position, c.Position.Add(c.Forward()), c.Up())
	}
}

func NewTargetCamera(target entity.Target, offset mgl32.Vec3) *TargetCamera {
	return &TargetCamera{
		Camera:       *NewCamera(),
		Target:       target,
		TargetOffset: offset,
	}
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
