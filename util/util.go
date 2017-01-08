package util

import (
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
)

// Why is this not in the standard time library? Am I missing something?
func GetTimeMillis() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func RandUint8() uint8 {
	return uint8(rand.Uint32() % 256)
}

func RandVec3() mgl32.Vec3 {
	return mgl32.Vec3{rand.Float32(), rand.Float32(), rand.Float32()}
}
