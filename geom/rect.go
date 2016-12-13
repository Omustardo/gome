package geom

import (
	"encoding/binary"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/bytecoder"
)

var (
	// Buffers are the float32 coordinates of two triangles (composing a 1x1 square), converted to a byte array, and
	// stored on the GPU. The gl.Buffer here is a reference to them.
	// This is the format required by OpenGL vertex buffers. This one buffer is used for all rectangles by modifying
	// the Scale, Rotation, and Translation matrices in the vertex shader.
	rectVertexBuffer gl.Buffer
	// Index buffer - rather than passing a minimum of 4 points (12 floats) to define a rectangle, just pass the indices
	// of those points in the rectTriangleBuffer.
	rectTrianglesIndexBuffer gl.Buffer
	rectLineStripIndexBuffer gl.Buffer

	// Texture coordinate buffer maps points in a rectangle to points on a texture. Note it assumes all textures match the
	// rectangle perfectly, so there will be distortion if you try to use a square texture with a rectangle.
	rectTextureCoordBuffer gl.Buffer
)

func initializeRect() {
	// Store basic rectangle vertices in a buffer.
	lower, upper := float32(-0.5), float32(0.5)
	rectVertices := bytecoder.Float32(binary.LittleEndian,
		lower, upper, 0,
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
	)
	rectVertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, rectVertexBuffer)             // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, rectVertices, gl.STATIC_DRAW) // store values in buffer

	// Store references to the vertices in different buffers.
	// For drawing full triangles, must specify two sets of 3 vertices. (gl.TRIANGLES)
	// Be careful to specify the correct order or the wrong side of the triangle will be in front and won't be rendered.
	rectTrianglesIndexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectTrianglesIndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, 1, 2, 3, 1, 3, 0), gl.STATIC_DRAW)
	// For drawing 4 line segments, must specify five points. (gl.LINE_LOOP)
	rectLineStripIndexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectLineStripIndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, 0, 1, 2, 3, 0), gl.STATIC_DRAW)

	rectTextureCoordBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, rectTextureCoordBuffer)
	textureCoordinates := bytecoder.Float32(binary.LittleEndian,
		0.0, 1.0,
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
	)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)
}

type Rect struct {
	entity.Entity
	model.Mesh
}

func (r *Rect) DrawWireframe() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(r.Color)
	shader.Basic.SetRotationMatrix2D(r.Rotation.Z())
	shader.Basic.SetScaleMatrix(r.Scale[0], r.Scale[1], r.Scale[2])
	shader.Basic.SetTranslationMatrix(r.Position.X(), r.Position.Y(), r.Position.Z())

	// Bind the array buffer before binding the element buffer, so it knows which array it's referencing.
	gl.BindBuffer(gl.ARRAY_BUFFER, rectVertexBuffer)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectLineStripIndexBuffer)
	itemSize := 3  // we use vertices made up of 3 floats
	itemCount := 5 // 4 segments, which requires 5 points
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.DrawElements(gl.LINE_STRIP, itemCount, gl.UNSIGNED_SHORT, 0)
	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}

func (r *Rect) DrawFilled() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(r.Color)
	shader.Basic.SetRotationMatrix2D(r.Rotation.Z())
	shader.Basic.SetScaleMatrix(r.Scale[0], r.Scale[1], r.Scale[2])
	shader.Basic.SetTranslationMatrix(r.Position.X(), r.Position.Y(), r.Position.Z())

	// Bind the array buffer before binding the element buffer, so it knows which array it's referencing.
	gl.BindBuffer(gl.ARRAY_BUFFER, rectVertexBuffer)
	// Bind element buffer so it is the target for DrawElements().
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectTrianglesIndexBuffer)

	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, 3 /* floats per vertex */, gl.FLOAT, false, 0, 0) // glVertexAttribPointer uses the buffer object that was bound to GL_ARRAY_BUFFER at the moment the function was called @@@ SUPER IMPORTANT
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib)                                               // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.DrawElements(gl.TRIANGLES, 6 /* num vertices for 2 triangles */, gl.UNSIGNED_SHORT, 0)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}

func (r *Rect) DrawTextured(texture gl.Texture) {
	shader.Texture.SetDefaults()
	shader.Texture.SetColor(r.Color)
	shader.Texture.SetRotationMatrix2D(r.Rotation.Z())
	shader.Texture.SetScaleMatrix(r.Scale[0], r.Scale[1], r.Scale[2])
	shader.Texture.SetTranslationMatrix(r.Position.X(), r.Position.Y(), r.Position.Z())
	shader.Texture.SetTextureSampler(texture)

	gl.BindBuffer(gl.ARRAY_BUFFER, rectTextureCoordBuffer)
	gl.VertexAttribPointer(shader.Texture.TextureCoordAttrib, 2, gl.FLOAT, false, 0, 0)
	gl.EnableVertexAttribArray(shader.Texture.TextureCoordAttrib)

	// Bind the array buffer before binding the element buffer, so it knows which array it's referencing.
	gl.BindBuffer(gl.ARRAY_BUFFER, rectVertexBuffer)
	// Bind element buffer so it is the target for DrawElements().
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectTrianglesIndexBuffer)

	gl.VertexAttribPointer(shader.Texture.VertexPositionAttrib, 3 /* floats per vertex */, gl.FLOAT, false, 0, 0) // glVertexAttribPointer uses the buffer object that was bound to GL_ARRAY_BUFFER at the moment the function was called @@@ SUPER IMPORTANT
	gl.EnableVertexAttribArray(shader.Texture.VertexPositionAttrib)                                               // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.DrawElements(gl.TRIANGLES, 6 /* num vertices for 2 triangles */, gl.UNSIGNED_SHORT, 0)

	gl.DisableVertexAttribArray(shader.Texture.TextureCoordAttrib)
	gl.DisableVertexAttribArray(shader.Texture.VertexPositionAttrib)
}
