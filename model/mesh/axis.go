package mesh

import (
	"encoding/binary"

	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
	"github.com/omustardo/gome/util"
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

	texture := util.LoadTextureData(3, 1, []uint8{
		// X Axis: Red
		255, 0, 0, 255,
		// Y Axis: Green
		0, 255, 0, 255,
		// Z Axis: Blue
		0, 0, 255, 255,
	})

	textureCoordBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
	textureCoordinates := bytecoder.Float32(binary.LittleEndian,
		// X Axis
		0, 0,
		0, 0,
		// Y Axis
		0.5, 0.5,
		0.5, 0.5,
		// Z Axis
		1, 1,
		1, 1,
	)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)

	return NewMesh(vbo, gl.Buffer{}, gl.Buffer{}, gl.LINES, 6, nil, texture, textureCoordBuffer)
}
