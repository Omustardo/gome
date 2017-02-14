package mesh

import (
	"encoding/binary"

	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
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
	texture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, 2, 2, gl.RGBA, gl.UNSIGNED_BYTE, textureData)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture{}) // bind to "null" to prevent using the wrong texture by mistake.

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
