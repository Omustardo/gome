// Package entity contains the Entity struct, which is the most basic physical world information attached to an object.
// Anything that lives in the game world should embed Entity.
package entity

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/util"
)

// Target provides a position. Functions that only deal with the position of an Entity should take a Target rather than
// a full entity.
type Target interface {
	GetPosition() mgl32.Vec3
}

type Entity struct {
	// Position of the entity.
	Position mgl32.Vec3

	// TODO: Add a center value which is a added to position when rendering. As it is, position can be thought of as the
	// bottom left corner of a cube that bounds a mesh. Being able to change positioning to an arbitrary center point will be necessary.

	// Rotation about the center. Note that this is a quaternion - a mathematical way of representing rotation.
	// Quaternions make some things easy, like having smooth rotations between different orientations, but
	// they aren't as intuitive as a simple Roll, Pitch, Yaw representation.
	//
	// If you don't want to deal with quaternions, just use the provided rotation related functions attached to the Entity struct:
	// SetRotation, ModifyRotationLocal, and ModifyRotationGlobal.
	//
	// When dealing with quaternions, there are a few basic things to know:
	// 1. Multiplying quaternions is like applying a rotation, and order matters.
	//    Q1.Mul(Q2) means start from the rotation of Q1, and apply the Q2 to it.
	// 2. You can convert roll, pitch, and yaw ("Euler Angles")  into a quaternion.
	//    quat := mgl32.AnglesToQuat(rot.X(), rot.Y(), rot.Z(), mgl32.XYZ)
	//    but converting back isn't built into mgl32. I believe the reason for this is there are multiple possible
	//    combinations of Euler angles that can all result in the same quaternion or overall orientation.
	// 3. A common situation is having an overall amount that you want to rotate per second, but you need to scale
	//    it down so only a little bit is done per frame. Do this using some form of quaternion interpolation.
	//    Slerp (Spherical Linear Interpolation) is quite common, but Nlerp (Normalized Linear Interpolation) is much faster
	//    and will likely give similar results.
	//    deltaRotation := mgl32.QuatSlerp(mgl32.QuatIdent(), rotationSpeedPerSecond, timeInSeconds)
	//
	//    There's also a function for this in the util package: util.ScaleQuatRotation(q mgl32.Quat, percent float32) mgl32.Quat
	//    Usage:
	//      rotationSpeed := mgl32.AnglesToQuat(0, 0, 2 * math.Pi, mgl32.XYZ) // Allow one full rotation per second.
	//      deltaTime := 0.016 																							  // Time passed in the last frame is very small.
	//      rotation := util.ScaleQuatRotation(rotationSpeed, deltaTime)		  // rotationSpeed * time = rotation
	//
	Rotation mgl32.Quat

	// Scale is how large the entity is in each dimension.
	// All meshes and collision boxes are expected to be initialized as unit cubes or unit squares centered at the origin
	// so scale can be used to compare them. This breaks if a mesh starts out as being larger/smaller since its scale
	// will need to be smaller/larger in order to draw it at the proper size. Any larger or smaller meshes should have their
	// meshes normalized to fit snugly in a unit cube as they are loaded.
	Scale mgl32.Vec3
}

func Default() Entity {
	return Entity{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.QuatIdent(), // TODO: If people forget to set Rotation, it prevents them from modifying the rotation later using .Rotate methods. Need to have a way of notifying developers of this.
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

func (e *Entity) Forward() mgl32.Vec3 {
	return e.Position.Add(e.Rotation.Rotate(mgl32.Vec3{1, 0, 0}))
}

func (e *Entity) Up() mgl32.Vec3 {
	return e.Rotation.Rotate(mgl32.Vec3{0, 1, 0})
}

// GetPosition gets the entity's position. Position is public and can be accessed directly, but this
// is necessary to make Entity implement the Target interface.
// TODO: Consider alternatives to this. We definitely want to be able to access just the position via an interface
// so functions that only need position don't get access to everything else.
// GetPosition isn't idiomatic Go, but the name Position would clash with the field.
// One option is to make the entity position private, but then rotation and scale should likely be made private
// and it makes everything just a bit more difficult to work with.
func (e *Entity) GetPosition() mgl32.Vec3 {
	return e.Position
}

// SetPosition directly sets the entity's position.
func (e *Entity) SetPosition(x, y, z float32) {
	e.Position[0] = validFloat32(x)
	e.Position[1] = validFloat32(y)
	e.Position[2] = validFloat32(z)
}

// ModifyPositionV adds the provided vector to the entity's position.
func (e *Entity) ModifyPositionV(vec mgl32.Vec3) {
	e.ModifyPosition(vec.X(), vec.Y(), vec.Z())
}

// ModifyPosition adds the provided vector to the entity's position.
func (e *Entity) ModifyPosition(x, y, z float32) {
	e.Position[0] += validFloat32(x)
	e.Position[1] += validFloat32(y)
	e.Position[2] += validFloat32(z)
}

// SetRotation takes a vector of Roll, Pitch, and Yaw in radians and sets the Entity to have that rotation.
// The rotations are applied in order. X then Y then Z (roll, then pitch, then yaw).
// Set rotation directly with `e.Rotation = mgl32.AnglesToQuat(x, y, z, mgl32.XYZ)` if you want the rotations to be
// applied in a different order.
func (e *Entity) SetRotation(rot mgl32.Vec3) {
	e.Rotation = mgl32.AnglesToQuat(rot.X(), rot.Y(), rot.Z(), mgl32.XYZ)
}

// ModifyRotationLocal applies the provided rotation based on the current orientation.
// The rotations are applied in order. X then Y then Z (roll, then pitch, then yaw).
func (e *Entity) ModifyRotationLocal(rot mgl32.Vec3) {
	e.Rotation = e.Rotation.Mul(mgl32.AnglesToQuat(rot.X(), rot.Y(), rot.Z(), mgl32.XYZ))
}

// ModifyRotationLocalQ applies the provided rotation based on the current orientation.
func (e *Entity) ModifyRotationLocalQ(rot mgl32.Quat) {
	e.Rotation = e.Rotation.Mul(rot)
}

// ModifyRotationGlobal applies the provided rotation based on the input being global / world space.
// This means the current orientation of the entity is ignored so it is no longer a relative rotation, but an absolute one.
// The rotations are applied in order. X then Y then Z (roll, then pitch, then yaw).
func (e *Entity) ModifyRotationGlobal(rot mgl32.Vec3) {
	e.ModifyRotationGlobalQ(mgl32.AnglesToQuat(rot.X(), rot.Y(), rot.Z(), mgl32.XYZ))
}

// ModifyRotationGlobal applies the provided rotation based on the input being global / world space.
// This means the current orientation of the entity is ignored so it is no longer a relative rotation, but an absolute one.
func (e *Entity) ModifyRotationGlobalQ(rot mgl32.Quat) {
	e.Rotation = rot.Mul(e.Rotation)
}

// RotationAngles returns the roll, pitch, and yaw of the Entity in radians. // TODO: confirm this works in the util package.
func (e *Entity) RotationAngles() mgl32.Vec3 {
	return util.QuatToEulerAngle(e.Rotation)
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
