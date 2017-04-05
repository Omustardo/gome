package camera

import (
	"log"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/core/entity"
)

var _ CameraI = (*Camera)(nil)

type CameraI interface {
	ModelView() mgl32.Mat4
	ProjectionOrthographic(width, height float32) mgl32.Mat4
	ProjectionPerspective(width, height float32) mgl32.Mat4
	Update(delta time.Duration)
	GetPosition() mgl32.Vec3
}

// Camera is the most basic camera type. It has no built in movement, zoom, or any special features.
// Its position and rotation can be manipulated via its embedded Entity. Make sure to also update the Up vector
// if you modify the rotation.
type Camera struct {
	// Entity makes the Camera a part of the world. Note that Entity.Scale is unused.
	// Entity.Up() is a vector pointing in the direction that the user will see as the top of the screen.
	entity.Entity

	// Near and Far are the range to render entities in front of the camera.
	// If the camera is expected to support zooming, be sure to make them small and large enough the change in position.
	// For example, if your object is 100 units away and your Zoomer can zoom in to 300%,
	// and zoom out to 25%, then the you must set Near<33.3 and Far>400 in order to always keep the target in view.
	Near, Far float32

	// Field of view in radians.
	// This only matters if using a perspective projection and can be ignored if using an orthographic projection.
	FOV float32
}

func NewCamera() *Camera {
	return &Camera{
		Entity: entity.Default(),
		Near:   0.1,
		Far:    10000,
		FOV:    math.Pi / 4.0,
	}
}

// ModelView returns a matrix used to transform from model space to camera coordinates.
// http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
func (c *Camera) ModelView() mgl32.Mat4 {
	if c.Up().ApproxEqual(mgl32.Vec3{0, 0, 0}) {
		log.Println("invalid ModelView: up vector is (0,0,0)")
	}
	if c.Forward().ApproxEqual(mgl32.Vec3{0, 0, 0}) {
		log.Println("invalid ModelView: forward vector is (0,0,0)")
	}
	return mgl32.LookAtV(c.Position, c.Position.Add(c.Forward()), c.Up())
}

// ProjectionOrthographic returns a matrix used to transform from camera space to screen space.
func (c *Camera) ProjectionOrthographic(width, height float32) mgl32.Mat4 {
	return mgl32.Ortho(-width/2, width/2,
		-height/2, height/2,
		c.Near, c.Far)
}

// ProjectionOrthographic returns a matrix used to transform from camera space to screen space.
// http://www.opengl-tutorial.org/beginners-tutorials/tutorial-3-matrices/
func (c *Camera) ProjectionPerspective(width, height float32) mgl32.Mat4 {
	return mgl32.Perspective(c.FOV, float32(width)/float32(height), c.Near, c.Far)
}

func (c *Camera) Update(delta time.Duration) {
	if c.Near >= c.Far {
		log.Println("camera's near is closer than far - nothing will render")
	}
}
