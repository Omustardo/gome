// bytecoder provides functions for converting standard arrays into byte arrays.
// Based on the Bytes function from golang.org\x\mobile\exp\f32\f32.go
package bytecoder

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// getByteOrder returns false for LittleEndian and true for BigEndian. It panics if the provided parameter isn't one of the two.
func getByteOrder(byteOrder binary.ByteOrder) bool {
	le := false
	switch byteOrder {
	case binary.BigEndian:
	case binary.LittleEndian:
		le = true
	default:
		panic(fmt.Sprintf("invalid byte order %v", byteOrder))
	}
	return le
}

// Bytes returns the byte representation of float32 values in the given byte
// order. byteOrder must be either binary.BigEndian or binary.LittleEndian.
func Float32(byteOrder binary.ByteOrder, values ...float32) []byte {
	le := getByteOrder(byteOrder)
	width := 4
	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		if le {
			b[width*i+0] = byte(u >> 0)
			b[width*i+1] = byte(u >> 8)
			b[width*i+2] = byte(u >> 16)
			b[width*i+3] = byte(u >> 24)
		} else {
			b[width*i+0] = byte(u >> 24)
			b[width*i+1] = byte(u >> 16)
			b[width*i+2] = byte(u >> 8)
			b[width*i+3] = byte(u >> 0)
		}
	}
	return b
}

func Vec2(byteOrder binary.ByteOrder, values ...mgl32.Vec2) []byte {
	le := getByteOrder(byteOrder)
	width := 8
	b := make([]byte, 8*len(values))
	for i, v := range values {
		x := math.Float32bits(v[0])
		y := math.Float32bits(v[1])
		if le {
			b[width*i+0] = byte(x >> 0)
			b[width*i+1] = byte(x >> 8)
			b[width*i+2] = byte(x >> 16)
			b[width*i+3] = byte(x >> 24)
			b[width*i+4] = byte(y >> 0)
			b[width*i+5] = byte(y >> 8)
			b[width*i+6] = byte(y >> 16)
			b[width*i+7] = byte(y >> 24)
		} else {
			b[width*i+0] = byte(x >> 24)
			b[width*i+1] = byte(x >> 16)
			b[width*i+2] = byte(x >> 8)
			b[width*i+3] = byte(x >> 0)
			b[width*i+4] = byte(y >> 24)
			b[width*i+5] = byte(y >> 16)
			b[width*i+6] = byte(y >> 8)
			b[width*i+7] = byte(y >> 0)
		}
	}
	return b
}

// Uint32 returns the byte representation a uint32 array in the given byte
// order. byteOrder must be either binary.BigEndian or binary.LittleEndian.
func Uint32(byteOrder binary.ByteOrder, values ...uint32) []byte {
	le := getByteOrder(byteOrder)
	width := 4
	b := make([]byte, 4*len(values))
	for i, v := range values {
		if le {
			b[width*i+0] = byte(v >> 0)
			b[width*i+1] = byte(v >> 8)
			b[width*i+2] = byte(v >> 16)
			b[width*i+3] = byte(v >> 24)
		} else {
			b[width*i+0] = byte(v >> 24)
			b[width*i+1] = byte(v >> 16)
			b[width*i+2] = byte(v >> 8)
			b[width*i+3] = byte(v >> 0)
		}
	}
	return b
}

// Uint16 returns the byte representation a uint16 array in the given byte
// order. byteOrder must be either binary.BigEndian or binary.LittleEndian.
func Uint16(byteOrder binary.ByteOrder, values ...uint16) []byte {
	le := getByteOrder(byteOrder)
	width := 2
	b := make([]byte, 2*len(values))
	for i, v := range values {
		if le {
			b[width*i+0] = byte(v >> 0)
			b[width*i+1] = byte(v >> 8)
		} else {
			b[width*i+0] = byte(v >> 8)
			b[width*i+1] = byte(v >> 0)
		}
	}
	return b
}
