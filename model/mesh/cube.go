package mesh

import (
	"image/color"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/glutil"
)

// Based on https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial/Creating_3D_objects_using_WebGL
func initializeCube() Mesh {
	// Store basic vertices in a buffer. This is a unit cube centered at the origin.
	lower, upper := float32(-0.5), float32(0.5)
	vertices := []float32{
		// Front
		lower, lower, upper,
		upper, lower, upper,
		upper, upper, upper,
		lower, upper, upper,

		// Back
		lower, lower, lower,
		lower, upper, lower,
		upper, upper, lower,
		upper, lower, lower,

		// Top
		lower, upper, lower,
		lower, upper, upper,
		upper, upper, upper,
		upper, upper, lower,

		// Bottom
		lower, lower, lower,
		upper, lower, lower,
		upper, lower, upper,
		lower, lower, upper,

		// Right
		upper, lower, lower,
		upper, upper, lower,
		upper, upper, upper,
		upper, lower, upper,

		// Left
		lower, lower, lower,
		lower, lower, upper,
		lower, upper, upper,
		lower, upper, lower,
	}

	indices := []uint16{
		0, 1, 2, 0, 2, 3, // front
		4, 5, 6, 4, 6, 7, // back
		8, 9, 10, 8, 10, 11, // top
		12, 13, 14, 12, 14, 15, // bottom
		16, 17, 18, 16, 18, 19, // right
		20, 21, 22, 20, 22, 23, // left
	}

	normals := []float32{
		// Front
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,

		// Back
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,

		// Top
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,

		// Bottom
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,

		// Right
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,

		// Left
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
	}

	textureCoordinates := []float32{
		// Front
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		// Back
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		// Top
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		// Bottom
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		// Right
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
		// Left
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
	}

	vertexBuffer := glutil.LoadBufferFloat32(vertices)
	indexBuffer := glutil.LoadIndexBuffer(indices)
	normalBuffer := glutil.LoadBufferFloat32(normals)
	textureCoordBuffer := glutil.LoadBufferFloat32(textureCoordinates)

	return NewMesh(vertexBuffer, indexBuffer, normalBuffer, gl.TRIANGLES, 36, nil, emptyTexture, textureCoordBuffer)
}

// NewCube returns a Mesh of a unit cube (all sides length 1) centered at the origin.
func NewCube(col *color.NRGBA, texture gl.Texture) Mesh {
	c := cube
	c.Color = col
	c.SetTexture(texture)
	return c
}
