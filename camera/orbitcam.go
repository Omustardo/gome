package camera

import (
	"time"

	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/input/keyboard"
)

var _ CameraI = (*OrbitCamera)(nil)

// OrbitCamera is a camera that always faces a target, but unlike TargetCamera, it moves in an orbit around the target.
type OrbitCamera struct {
	Camera

	// Target is the entity which the OrbitCamera always faces. The OrbitCamera.Camera.Entity.Rotation determines
	// which direction it faces it from.
	Target entity.Target
	// TargetOffset determines how far the camera is positioned from the target.
	TargetOffset float32

	// RotateSpeed determines how many radians in total the camera can rotate per second.
	RotateSpeed float32

	// Zoomer handles camera zoom.
	Zoomer zoom.Zoom
}

func (c *OrbitCamera) Update(delta time.Duration) {
	c.Camera.Update(delta)

	rotate := mgl32.Vec3{0, 0, 0}
	// TODO: Move this logic so keybinds are customizeable and so camera doesn't depend on keyboard package.
	if keyboard.Handler.IsKeyDown(glfw.KeyLeft) {
		rotate[1] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyRight) {
		rotate[1] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyUp) {
		rotate[2] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyDown) {
		rotate[2] = 1
	}
	if rotate.Len() > 0 {
		c.ModifyRotationLocal(rotate.Normalize().Mul(float32(delta.Seconds()) * c.RotateSpeed)) // @@@@@@@@ does global vs local matter here?
	}

	// Take the Camera's rotation and desired distance from target. Negate it, and add it to the target position
	// to move the camera "back" to where it should be.
	// For example, if the camera is looking down the X axis and the offset is 5 units, we need to end up adding {-5,0,0}
	// to the target's position to get to the proper location.
	offset := c.Entity.Rotation.Rotate(mgl32.Vec3{-c.TargetOffset / c.GetCurrentZoomPercent(), 0, 0})
	c.Entity.Position = c.Target.GetPosition().Add(offset)
}

func NewOrbitCamera(target entity.Target, offset float32) *OrbitCamera {
	return &OrbitCamera{
		Camera:       *NewCamera(),
		Target:       target,
		TargetOffset: offset,
		RotateSpeed:  math.Pi / 2,
	}
}

func (c *OrbitCamera) GetCurrentZoomPercent() float32 {
	return 1
}
