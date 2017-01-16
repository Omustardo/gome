// Package entity contains the Entity struct, which is the most basic physical world information attached to an object.
// Anything that lives in the game world should embed Entity.
package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/util"
)

func Default() Entity {
	return Entity{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.QuatIdent(), //  {W:0, V:mgl32.Vec3{0,0,0}} @@@@@@ TODO: If people forget to set Rotation, it prevents them from modifying the rotation later using .Rotate methods.
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

type Entity struct {
	// Center coordinates of the entity.
	Position mgl32.Vec3

	// Rotation about the center in radians.
	Rotation mgl32.Quat

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

func (e *Entity) ModifyCenterV(vec mgl32.Vec3) {
	e.ModifyCenter(vec.X(), vec.Y(), vec.Z())
}
func (e *Entity) ModifyCenter(x, y, z float32) {
	e.Position[0] += validFloat32(x)
	e.Position[1] += validFloat32(y)
	e.Position[2] += validFloat32(z)
}

// SetRotation takes a vector of angles in radians and sets the Entity to have that rotation.
// The rotations are applied in order. X then Y then Z.
// Set rotation directly with e.Rotation = mgl32.AnglesToQuat(x, y, z, mgl32.XYZ) if you want the rotations to be
// applied in a different order.
func (e *Entity) SetRotationV(rot mgl32.Vec3) {
	e.SetRotation(rot.X(), rot.Y(), rot.Z())
}

// SetRotation takes x,y,z angles in radians and sets the Entity to have that rotation.
// The rotations are applied in order. X then Y then Z.
// Set rotation directly with e.Rotation = mgl32.AnglesToQuat(x, y, z, mgl32.XYZ) if you want the rotations to be
// applied in a different order.
func (e *Entity) SetRotation(x, y, z float32) {
	e.Rotation = mgl32.AnglesToQuat(x, y, z, mgl32.XYZ)
}

// ModifyRotationGlobal applies the provided rotation based on the input being global / world space.
func (e *Entity) ModifyRotationGlobal(rot mgl32.Vec3) {
	e.ModifyRotationGlobalQ(mgl32.AnglesToQuat(rot.X(), rot.Y(), rot.Z(), mgl32.XYZ))
}

// ModifyRotationGlobal applies the provided rotation based on the input being global / world space.
func (e *Entity) ModifyRotationGlobalQ(rot mgl32.Quat) {
	e.Rotation = rot.Mul(e.Rotation)
}

// ModifyRotation applies the provided rotation based on the current orientation.
func (e *Entity) ModifyRotationLocalQ(rot mgl32.Quat) {
	e.Rotation = e.Rotation.Mul(rot)
}

// RotationAngles returns the rotation of the Entity in radians.
func (e *Entity) RotationAngles() mgl32.Vec3 {
	return util.QuatToEulerAngle(e.Rotation)
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
