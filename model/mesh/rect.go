package mesh

import (
	"encoding/binary"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

func NewRect(color *color.NRGBA, texture gl.Texture) Mesh {
	r := rect
	r.Color = color
	r.texture = texture
	SetValidDefaults(&r)
	return r
}

func NewRectOutline(color *color.NRGBA) Mesh {
	r := wireframeRect
	r.Color = color
	SetValidDefaults(&r)
	return r
}

func initializeRect() Mesh {
	// Store basic rectangle vertices in a buffer.
	lower, upper := float32(-0.5), float32(0.5)
	vertices := bytecoder.Float32(binary.LittleEndian,
		lower, upper, 0,
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
	)
	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)                // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW) // store values in buffer
	indexBuffer := gl.CreateBuffer()

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, 0, 1, 2, 0, 2, 3), gl.STATIC_DRAW)

	textureCoordBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
	textureCoordinates := bytecoder.Float32(binary.LittleEndian,
		0.0, 1.0,
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
	)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)

	// Normals for a unit square simple extend in the directions of the corners.
	unitLength := mgl32.Vec2{1, 1}.Normalize().Len()
	lower, upper = -unitLength, unitLength
	normals := bytecoder.Float32(binary.LittleEndian,
		lower, upper, 0,
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
	)
	normalVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, normalVBO)
	gl.BufferData(gl.ARRAY_BUFFER, normals, gl.STATIC_DRAW)

	return NewMesh(vertexVBO, indexBuffer, normalVBO, gl.TRIANGLES, 6, nil, gl.Texture{}, textureCoordBuffer)
}

func initializeWireframeRect() Mesh {
	lower, upper := float32(-0.5), float32(0.5)
	rectVertices := bytecoder.Float32(binary.LittleEndian,
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
		lower, upper, 0,
	)
	vbo := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)                          // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, rectVertices, gl.STATIC_DRAW) // store values in buffer

	return NewMesh(vbo, gl.Buffer{}, gl.Buffer{}, gl.LINE_LOOP, 4, nil, gl.Texture{}, gl.Buffer{})
}
