package camera

import (
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome/input/keyboard"
)

var _ CameraI = (*FreeCamera)(nil)

// FreeCamera is a camera that moves freely in space, relative to its own orientation.
type FreeCamera struct {
	// Camera is an embedded basic camera struct. Note that Camera.Up is ignored and a value of (0,1,0) is used.
	// This Up vector is rotated by the Camera.Entity.Rotation, so if you want a different initial rotation then
	// modify that. Otherwise leave it to the user to modify rotation and position via the Update function.
	Camera

	// MoveSpeed determines how much the camera's position can change per second.
	MoveSpeed float32

	// RotateSpeed determines how many radians in total the camera can rotate per second.
	RotateSpeed float32
}

func NewFreeCamera() *FreeCamera {
	return &FreeCamera{
		Camera:      *NewCamera(),
		MoveSpeed:   100,
		RotateSpeed: 2 * math.Pi / 4,
	}
}

func (c *FreeCamera) Update(delta time.Duration) {
	c.Camera.Update(delta)

	move := mgl32.Vec3{}
	rotate := mgl32.Vec3{}
	// WASD to move forward, back, left, right. Q,E for down and up.
	// Arrows for rotation.
	// TODO: Move this logic so keybinds are customizeable and so camera doesn't depend on keyboard package.
	if keyboard.Handler.IsKeyDown(glfw.KeyW) {
		move[0] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyS) {
		move[0] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyQ) {
		move[1] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyE) {
		move[1] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyA) {
		move[2] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyD) {
		move[2] = 1
	}

	if keyboard.Handler.IsKeyDown(glfw.KeyLeft) {
		rotate[1] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyRight) {
		rotate[1] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyUp) {
		rotate[2] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyDown) {
		rotate[2] = -1
	}
	// TODO: I'm unsure whether to apply rotation or position changes first. The outcome is different, but since this
	// should be done in very small steps over multiple frames, it shouldn't matter from a human perspective.
	if move.Len() > 0 {
		move = c.Rotation.Rotate(move)
		c.ModifyPositionV(move.Normalize().Mul(float32(delta.Seconds()) * c.MoveSpeed))
	}
	if rotate.Len() > 0 {
		c.ModifyRotationLocal(rotate.Normalize().Mul(float32(delta.Seconds()) * c.RotateSpeed))
	}
}
