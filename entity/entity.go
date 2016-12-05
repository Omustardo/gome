package entity

import "github.com/go-gl/mathgl/mgl32"

type Entity struct {
	// Center coordinates of the entity.
	Position mgl32.Vec3

	// Rotation is in radians.
	Rotation mgl32.Vec3

	Scale mgl32.Vec3
}

func (e *Entity) Center() mgl32.Vec3 {
	return e.Position
}

// Target has a center. If something only needs access to an entity's position, pass it as a Target.
type Target interface {
	Center() mgl32.Vec3
}
