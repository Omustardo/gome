package util

import (
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// TODO: Why is this not in the standard time library? Am I missing something?
func GetTimeMillis() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func RandUint8() uint8 {
	return uint8(rand.Uint32() % 256)
}

func RandVec3() mgl32.Vec3 {
	return mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}.Normalize()
}

func RandQuat() mgl32.Quat {
	return mgl32.AnglesToQuat(rand.Float32(), rand.Float32(), rand.Float32(), mgl32.ZYX).Normalize() // TODO: Do I need to normalize?
}

// ScaleQuatRotation scales the provided rotation by the provided percentage.
// For example, if you have a quaternion representing a 90 degree rotation: mgl32.AnglesToQuat(mgl32.DegToRad(90), 0, 0, mgl32.ZYX)
// you could get a quaternion that rotates half as much with: ScaleQuatRotation(q, 0.5)
// Note that percentages above 1.0 will increase the provided rotation as expected, and negative percentages will reverse the rotation.
//
// Sample usage to limit rotation based on a maximum amount per second:
//   rotationSpeed := mgl32.AnglesToQuat(0, 0, 2 * math.Pi, mgl32.XYZ)	// Allow one full rotation per second.
//   deltaTime := 0.016 												// Time passed in the last frame is very small.
//   rotation := util.ScaleQuatRotation(rotationSpeed, deltaTime)		// rotationSpeed * time = rotation
func ScaleQuatRotation(q mgl32.Quat, percent float32) mgl32.Quat {
	// TODO: Consider using mgl32.QuatNlerp() as it's much faster and likely won't look different for small rotations.
	return mgl32.QuatSlerp(mgl32.QuatIdent(), q, percent) // TODO: Confirm that this doesn't need to be normalized.
}

func IsPowerOfTwo(n int) bool {
	// http://www.graphics.stanford.edu/~seander/bithacks.html#DetermineIfPowerOf2
	return (n > 0) && (n&(n-1)) == 0
}

// RoundUpToPowerOfTwo returns the smallest number that is >= n and also a power of two.
// If n is close to the max int value, unexpected behavior results are likely.
func RoundUpToPowerOfTwo(n int) int {
	// http://stackoverflow.com/questions/466204/rounding-up-to-nearest-power-of-2
	power := 1
	for power < n {
		power *= 2
	}
	return power
}
