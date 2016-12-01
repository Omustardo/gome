package cube

import (
	"encoding/binary"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/bytecoder"
)

var (
	// Buffers are the float32 coordinates of two triangles (composing a 1x1 square), converted to a byte array, and
	// stored on the GPU. The gl.Buffer here is a reference to them.
	// This is the format required by OpenGL vertex buffers. This one buffer is used for all rectangles by modifying
	// the Scale, Rotation, and Translation matrices in the vertex shader.
	vertexBuffer gl.Buffer
	// Index buffer - rather than passing a minimum of 4 points (12 floats) to define a rectangle, just pass the indices
	// of those points in the rectTriangleBuffer.
	indexBuffer gl.Buffer

	// Texture coordinate buffer maps points in a rectangle to points on a texture. Note it assumes all textures match the
	// rectangle perfectly, so there will be distortion if you try to use a square texture with a rectangle.
	textureCoordBuffer gl.Buffer
)

func Initialize() {
	// Store basic vertices in a buffer.
	lower, upper := float32(-0.5), float32(0.5)
	rectVertices := bytecoder.Float32(binary.LittleEndian,
		// Front face
		lower, lower, upper,
		upper, lower, upper,
		upper, upper, upper,
		lower, upper, upper,

		// Back face
		lower, lower, lower,
		lower, upper, lower,
		upper, upper, lower,
		upper, lower, lower,

		// Top face
		lower, upper, lower,
		lower, upper, upper,
		upper, upper, upper,
		upper, upper, lower,

		// Bottom face
		lower, lower, lower,
		upper, lower, lower,
		upper, lower, upper,
		lower, lower, upper,

		// Right face
		upper, lower, lower,
		upper, upper, lower,
		upper, upper, upper,
		upper, lower, upper,

		// Left face
		lower, lower, lower,
		lower, lower, upper,
		lower, upper, upper,
		lower, upper, lower,
	)
	vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)                 // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, rectVertices, gl.STATIC_DRAW) // store values in buffer

	// Store references to the vertices in different buffers.
	// For drawing full triangles, must specify two sets of 3 vertices. (gl.TRIANGLES)
	// Be careful to specify the correct order or the wrong side of the triangle will be in front and won't be rendered.
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
	// Center coordinates of the cube.
	Center mgl32.Vec3
	Dim    mgl32.Vec3
	// Angle is radians of rotation around the center.
	Rotation   mgl32.Vec3
	R, G, B, A float32
}

func (c *Cube) DrawTextured(texture gl.Texture) {
	shader.Texture.SetDefaults()
	shader.Texture.SetRotationMatrix(c.Rotation.X(), c.Rotation.Y(), c.Rotation.Z())
	shader.Texture.SetScaleMatrix(c.Dim.X(), c.Dim.Y(), c.Dim.Z())
	shader.Texture.SetTranslationMatrix(c.Center.X(), c.Center.Y(), c.Center.Z())
	shader.Texture.SetTextureSampler(texture)

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
