package mesh

import (
	"encoding/binary"
	"image/color"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

// Based on https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial/Creating_3D_objects_using_WebGL
func initializeCube() Mesh {
	// Store basic vertices in a buffer. This is a unit cube centered at the origin.
	lower, upper := float32(-0.5), float32(0.5)
	vertices := bytecoder.Float32(binary.LittleEndian,
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
	)
	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)                // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW) // store values in buffer

	indexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian,
		0, 1, 2, 0, 2, 3, // front
		4, 5, 6, 4, 6, 7, // back
		8, 9, 10, 8, 10, 11, // top
		12, 13, 14, 12, 14, 15, // bottom
		16, 17, 18, 16, 18, 19, // right
		20, 21, 22, 20, 22, 23, // left
	), gl.STATIC_DRAW)

	normals := bytecoder.Float32(binary.LittleEndian,
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
	)
	normalVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, normalVBO)               // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, normals, gl.STATIC_DRAW) // store values in buffer

	textureCoordBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
	textureCoordinates := bytecoder.Float32(binary.LittleEndian,
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
	)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)

	return NewMesh(vertexVBO, indexBuffer, normalVBO, gl.TRIANGLES, 36, nil, EmptyTexture, textureCoordBuffer)
}

func NewCube(col *color.NRGBA, texture gl.Texture) Mesh {
	c := cube
	c.Color = col
	c.texture = texture
	SetValidDefaults(&c)
	return c
}
