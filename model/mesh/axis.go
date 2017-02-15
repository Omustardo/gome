package mesh

import (
	"encoding/binary"

	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
	"github.com/omustardo/gome/util/glutil"
)

func initializeAxes() Mesh {
	axisLength := float32(0.5)
	verts := bytecoder.Float32(binary.LittleEndian,
		-axisLength, 0, 0,
		axisLength, 0, 0,
		0, -axisLength, 0,
		0, axisLength, 0,
		0, 0, -axisLength,
		0, 0, axisLength,
	)
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, verts, gl.STATIC_DRAW)

	// TODO: Texture loading is done in the asset package, but using asset.LoadTexture here creates an odd dependency.
	// Consider making a separate package solely for putting bytes into GPU buffers.
	textureData := []uint8{
		// X Axis: Red
		255, 0, 0, 255,
		// Y Axis: Green
		0, 255, 0, 255,
		// Z Axis: Blue
		0, 0, 255, 255,
		// Pad so the texture size is a power of 2.
		0, 0, 0, 0,
	}
	texture, err := glutil.LoadTextureData(2, 2, textureData)
	if err != nil {
		panic(err) // TODO: Ideally move away from this sort of initialize / singleton method, but if not at least return errors.
	}

	textureCoordBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
	textureCoordinates := bytecoder.Float32(binary.LittleEndian,
		// X Axis
		0, 0,
		0, 0,
		// Y Axis
		1, 0,
		1, 0,
		// Z Axis
		0, 1,
		0, 1,
	)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)

	return NewMesh(vbo, gl.Buffer{}, gl.Buffer{}, gl.LINES, 6, nil, texture, textureCoordBuffer)
}
