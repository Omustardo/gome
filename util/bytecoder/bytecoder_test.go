package bytecoder

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestFloat32(t *testing.T) {
	// Tests should match IEEE754. Use https://www.h-schmidt.net/FloatConverter/IEEE754.html
	tests := []struct {
		array     []float32
		byteOrder binary.ByteOrder
		want      []byte
	}{
		{
			array:     []float32{1.0},
			byteOrder: binary.LittleEndian,
			want:      []byte{0, 0, 0x80, 0x3f},
		},
		{
			array:     []float32{1.0},
			byteOrder: binary.BigEndian,
			want:      []byte{0x3f, 0x80, 0, 0},
		},
		{
			array:     []float32{1.0, 16},
			byteOrder: binary.LittleEndian,
			want: []byte{
				0, 0, 0x80, 0x3f, // 1.0
				0, 0, 0x80, 0x41, // 16.0
			},
		},
		{
			array:     []float32{1.0, 16},
			byteOrder: binary.BigEndian,
			want: []byte{
				0x3f, 0x80, 0, 0, // 1.0
				0x41, 0x80, 0, 0, // 16.0
			},
		},
	}

	for _, tt := range tests {
		got := Float32(tt.byteOrder, tt.array...)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Converting: %v, Got %+v, wanted %+v", tt.array, got, tt.want)
		}
	}
}

func TestVec2(t *testing.T) {
	tests := []struct {
		array     []mgl32.Vec2
		byteOrder binary.ByteOrder
		want      []byte
	}{
		{
			array:     []mgl32.Vec2{{1.0, 16.0}},
			byteOrder: binary.LittleEndian,
			want: []byte{
				0, 0, 0x80, 0x3f, // 1.0
				0, 0, 0x80, 0x41, // 16.0
			},
		},
		{
			array:     []mgl32.Vec2{{1.0, 16.0}},
			byteOrder: binary.BigEndian,
			want: []byte{
				0x3f, 0x80, 0, 0, // 1.0
				0x41, 0x80, 0, 0, // 16.0
			},
		},
	}

	for _, tt := range tests {
		got := Vec2(tt.byteOrder, tt.array...)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Converting: %v, Got %+v, wanted %+v", tt.array, got, tt.want)
		}
	}
}

func TestUint16(t *testing.T) {
	tests := []struct {
		array     []uint16
		byteOrder binary.ByteOrder
		want      []byte
	}{
		{
			array:     []uint16{1, 2, 4, 8},
			byteOrder: binary.LittleEndian,
			want: []byte{
				1, 0,
				2, 0,
				4, 0,
				8, 0,
			},
		},
		{
			array:     []uint16{16, 32},
			byteOrder: binary.BigEndian,
			want: []byte{
				0, 16,
				0, 32,
			},
		},
		{
			array:     []uint16{65535, 65534},
			byteOrder: binary.LittleEndian,
			want: []byte{
				0xFF, 0xFF,
				0xFE, 0xFF,
			},
		},
	}

	for _, tt := range tests {
		got := Uint16(tt.byteOrder, tt.array...)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Converting: %v, Got %+v, wanted %+v", tt.array, got, tt.want)
		}
	}
}

func TestUint32(t *testing.T) {
	tests := []struct {
		array     []uint32
		byteOrder binary.ByteOrder
		want      []byte
	}{
		{
			array:     []uint32{1, 2, 4, 8},
			byteOrder: binary.LittleEndian,
			want: []byte{
				1, 0, 0, 0,
				2, 0, 0, 0,
				4, 0, 0, 0,
				8, 0, 0, 0,
			},
		},
		{
			array:     []uint32{16, 32},
			byteOrder: binary.BigEndian,
			want: []byte{
				0, 0, 0, 16,
				0, 0, 0, 32,
			},
		},
		{
			array:     []uint32{65535, 65534},
			byteOrder: binary.LittleEndian,
			want: []byte{
				0xFF, 0xFF, 0, 0,
				0xFE, 0xFF, 0, 0,
			},
		},
		{
			array:     []uint32{587578374, 3645234},
			byteOrder: binary.LittleEndian,
			want: []byte{
				0x06, 0xBC, 0x05, 0x23,
				0x32, 0x9F, 0x37, 0,
			},
		},
	}

	for _, tt := range tests {
		got := Uint32(tt.byteOrder, tt.array...)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Converting: %v, Got %+v, wanted %+v", tt.array, got, tt.want)
		}
	}
}
