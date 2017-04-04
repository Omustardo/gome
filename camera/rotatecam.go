package camera

import (
	"time"

	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/input/keyboard"
)

var _ CameraI = (*RotateCamera)(nil)

// RotateCamera is a camera that always faces a target. Unlike TargetCamera which always stays at a fixed position in
// relation to the target, and unlike OrbitCamera which moves in an orbit around the target, the RotateCamera rotates
// around a particular axis at the target's location.
// You can think of its movement like latitude and longitude on the Earth, if the Earth were the target.
type RotateCamera struct {
	Camera

	// Target is the entity which the OrbitCamera always faces. The OrbitCamera.Camera.Entity.Rotation determines
	// which direction it faces it from.
	Target entity.Target
	// TargetOffset determines how far the camera is positioned from the target.
	TargetOffset float32

	// RotateSpeed determines how many radians in total the camera can rotate per second.
	RotateSpeed float32
}

func (c *RotateCamera) Update(delta time.Duration) {
	c.Camera.Update(delta)

	var rotate mgl32.Vec3
	// TODO: Move this logic so keybinds are customizeable and so camera doesn't depend on keyboard package.
	if keyboard.Handler.IsKeyDown(glfw.KeyLeft) {
		rotate[1] += -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyRight) {
		rotate[1] += 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyUp) {
		rotate[0] += -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyDown) {
		rotate[0] += 1
	}
	if rotate.Len() != 0 {
		c.ModifyRotationGlobal(rotate.Normalize().Mul(float32(delta.Seconds()) * c.RotateSpeed))
	}
	// TODO: Cap rotation to a set number of degrees so we can't flip upside down and make the controls become backwards.

	// Take the Camera's rotation and desired distance from target. Negate it, and add it to the target position
	// to move the camera "back" to where it should be.
	// For example, if the camera is looking down the X axis and the offset is 5 units, we need to end up adding {-5,0,0}
	// to the target's position to get to the proper location.
	offset := c.Entity.Rotation.Rotate(entity.Forward.Mul(-c.TargetOffset))
	c.Entity.Position = c.Target.GetPosition().Add(offset)
}

func NewRotateCamera(target entity.Target, offset float32) *RotateCamera {
	return &RotateCamera{
		Camera:       *NewCamera(),
		Target:       target,
		TargetOffset: offset,
		RotateSpeed:  math.Pi / 2,
	}
}
