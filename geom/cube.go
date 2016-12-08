package geom

import (
	"encoding/binary"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/drawable"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/bytecoder"
)

var (
	vertexBuffer       gl.Buffer
	indexBuffer        gl.Buffer
	textureCoordBuffer gl.Buffer
)

func initializeCube() {
	// Store basic vertices in a buffer. This is a unit cube centered at the origin.
	lower, upper := float32(-0.5), float32(0.5)
	rectVertices := bytecoder.Float32(binary.LittleEndian,
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
	vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)                 // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, rectVertices, gl.STATIC_DRAW) // store values in buffer

	// Store references to the vertices in different buffers.
	// These are for drawing full triangles using mode gl.TRIANGLES
	indexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian,
		0, 1, 2, 0, 2, 3, // front
		4, 5, 6, 4, 6, 7, // back
		8, 9, 10, 8, 10, 11, // top
		12, 13, 14, 12, 14, 15, // bottom
		16, 17, 18, 16, 18, 19, // right
		20, 21, 22, 20, 22, 23, // left
	), gl.STATIC_DRAW)

	textureCoordBuffer = gl.CreateBuffer()
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
}

type Cube struct {
	entity.Entity
	drawable.Attributes
}

func (c *Cube) Draw() {
	shader.Texture.SetDefaults()
	shader.Texture.SetRotationMatrix(c.Rotation.X(), c.Rotation.Y(), c.Rotation.Z())
	shader.Texture.SetScaleMatrix(c.Scale.X(), c.Scale.Y(), c.Scale.Z())
	shader.Texture.SetTranslationMatrix(c.Position.X(), c.Position.Y(), c.Position.Z())
	shader.Texture.SetTextureSampler(*c.Texture)

	gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
	gl.VertexAttribPointer(shader.Texture.TextureCoordAttrib, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Texture.TextureCoordAttrib)

	// Bind the array buffer before binding the element buffer, so it knows which array it's referencing.
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	// Bind element buffer so it is the target for DrawElements().
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)

	gl.VertexAttribPointer(shader.Texture.VertexPositionAttrib, 3 /* floats per vertex */, gl.FLOAT, false, 0, 0) // glVertexAttribPointer uses the buffer object that was bound to GL_ARRAY_BUFFER at the moment the function was called @@@ SUPER IMPORTANT
	gl.EnableVertexAttribArray(shader.Texture.VertexPositionAttrib)                                               // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.DrawElements(gl.TRIANGLES, 6*6 /* num vertices for 6 triangles */, gl.UNSIGNED_SHORT, 0)

	gl.DisableVertexAttribArray(shader.Texture.TextureCoordAttrib)
	gl.DisableVertexAttribArray(shader.Texture.VertexPositionAttrib)
}
