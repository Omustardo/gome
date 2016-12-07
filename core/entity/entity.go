// Package entity contains the Entity struct, which is the most basic physical world information attached to an object.
// Anything that lives in the game world should embed Entity.
package entity

import "github.com/go-gl/mathgl/mgl32"

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

// Target has a center. If something only needs access to an entity's position, pass it as a Target.
type Target interface {
	Center() mgl32.Vec3
}
