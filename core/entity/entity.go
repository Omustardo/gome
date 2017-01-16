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

// RotationAngles returns the roll, pitch, and yaw of the Entity in radians. The
// Note that multiple
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
