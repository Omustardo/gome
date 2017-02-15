package mesh

import (
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/glutil"
)

func initializeAxes() Mesh {
	axisLength := float32(0.5)
	verts := []float32{
		-axisLength, 0, 0,
		axisLength, 0, 0,
		0, -axisLength, 0,
		0, axisLength, 0,
		0, 0, -axisLength,
		0, 0, axisLength,
	}
	vbo := glutil.LoadBufferFloat32(verts)

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

	textureCoordinates := []float32{
		// X Axis
		0, 0,
		0, 0,
		// Y Axis
		1, 0,
		1, 0,
		// Z Axis
		0, 1,
		0, 1,
	}
	textureCoordBuffer := glutil.LoadBufferFloat32(textureCoordinates)

	return NewMesh(vbo, gl.Buffer{}, gl.Buffer{}, gl.LINES, 6, nil, texture, textureCoordBuffer)
}
