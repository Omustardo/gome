// Package entity contains the Entity struct, which is the most basic physical world information attached to an object.
// Anything that lives in the game world should embed Entity.
package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

var Default = Entity{
	Position: mgl32.Vec3{0, 0, 0},
	Scale:    mgl32.Vec3{1, 1, 1},
	Rotation: mgl32.Vec3{0, 0, 0},
}

type Entity struct {
	// Center coordinates of the entity.
	Position mgl32.Vec3

	// Rotation about the center in radians.
	Rotation mgl32.Vec3

	// All meshes and collision boxes are expected to be initialized as unit cubes or unit squares centered at the origin.
	// Any other sized objects should be made that way by scaling up or down.
	Scale mgl32.Vec3
}

func (e *Entity) Center() mgl32.Vec3 {
	return e.Position
}

func (e *Entity) SetCenter(x, y, z float32) {
	e.Position[0] = validFloat32(x)
	e.Position[1] = validFloat32(y)
	e.Position[2] = validFloat32(z)
}

func (e *Entity) ModifyCenter(x, y, z float32) {
	e.Position[0] += validFloat32(x)
	e.Position[1] += validFloat32(y)
	e.Position[2] += validFloat32(z)
}

// Target has a center. If something only needs access to an entity's position, pass it as a Target.
type Target interface {
	Center() mgl32.Vec3
}

// validFloat32 checks for values that wouldn't exist in the coordinate system of the game world and returns
// zero in those cases, or the input if there isn't an issue.
// This could easily occur in practice when normalizing a vector. If the vector is length zero, then normalizing it
// makes it a vector of NaN values.
func validFloat32(x float32) float32 {
	x64 := float64(x)
	if math.IsNaN(x64) || math.IsInf(x64, 0) {
		return 0
	}
	return x
}
