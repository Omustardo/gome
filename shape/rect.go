package shape

import (
	"encoding/binary"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/bytecoder"
	"github.com/omustardo/gome/shader"
)

var _ Shape = (*Rect)(nil)

var (
	// Buffers are the float32 coordinates of two triangles (composing a 1x1 square), converted to a byte array, and
	// stored on the GPU. The gl.Buffer here is a reference to them.
	// This is the format required by OpenGL vertex buffers. This one buffer is used for all rectangles by modifying
	// the Scale, Rotation, and Translation matrices in the vertex shader.
	rectVertexBuffer gl.Buffer
	// Index buffer - rather than passing a minimum of 4 points (12 floats) to define a rectangle, just pass the indices
	// of those points in the rectTriangleBuffer.
	rectTriangleIndexBuffer  gl.Buffer
	rectLineStripIndexBuffer gl.Buffer
)

func loadRectangles() {
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
	rectTriangleIndexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectTriangleIndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, 1, 2, 3, 1, 3, 0), gl.STATIC_DRAW)
	// For drawing 4 line segments, must specify five points. (gl.LINE_LOOP)
	rectLineStripIndexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectLineStripIndexBuffer)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, 0, 1, 2, 3, 0), gl.STATIC_DRAW) // TODO: test leaving out the last index - does it automatically connect back? probably not.
}

type Rect struct {
	// X, Y are the center coordinate of the rectangle.
	X, Y          float32
	Width, Height float32
	// Angle is radians of rotation around the center.
	Angle      float32
	R, G, B, A float32
}

func (r *Rect) Draw() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(r.R, r.G, r.B, r.A)
	shader.Basic.SetRotationMatrix2D(r.Angle)
	shader.Basic.SetScaleMatrix(r.Width, r.Height, 0)
	shader.Basic.SetTranslationMatrix(r.X, r.Y, 0)

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
	shader.Basic.SetColor(r.R, r.G, r.B, r.A)
	shader.Basic.SetRotationMatrix2D(r.Angle)
	shader.Basic.SetScaleMatrix(r.Width, r.Height, 0)
	shader.Basic.SetTranslationMatrix(r.X, r.Y, 0)

	// Bind the array buffer before binding the element buffer, so it knows which array it's referencing.
	gl.BindBuffer(gl.ARRAY_BUFFER, rectVertexBuffer)
	// Bind element buffer so it is the target for DrawElements().
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, rectTriangleIndexBuffer)

	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, 3 /* floats per vertex */, gl.FLOAT, false, 0, 0) // glVertexAttribPointer uses the buffer object that was bound to GL_ARRAY_BUFFER at the moment the function was called @@@ SUPER IMPORTANT
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib)                                               // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.DrawElements(gl.TRIANGLES, 6 /* num vertices for 2 triangles */, gl.UNSIGNED_SHORT, 0)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}

//func DrawRectsFilled(rects []Rect) {
//	shader.Basic.SetDefaults()
//	//shader.SetRotationMatrix2D(rects[0].Angle)
//	//shader.SetScaleMatrix(rects[0].Width, rects[0].Height, 0)
//	//shader.SetTranslationMatrix(rects[0].X, rects[0].Y, 0)
//
//	indices := []uint16{}
//	for range rects {
//		// Every rectangle has the same starting vertices. They will be translated in the shader. // TODO: recreating this slice each time is terrible.
//		indices = append(indices, 1, 2, 3, 1, 3, 0)
//	}
//	indexBuffer := gl.CreateBuffer()
//	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, indexBuffer)
//	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, indices...), gl.STATIC_DRAW)
//
//	gl.BindBuffer(gl.ARRAY_BUFFER, rectVertexBuffer)
//
//	itemSize := 3                                                 // because the points consist of 3 floats
//	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
//	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)
//	gl.DrawElements(gl.TRIANGLES, len(indices), gl.UNSIGNED_SHORT, 0)
//
//	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
//	gl.DeleteBuffer(indexBuffer)
//}

func (r *Rect) SetCenter(x, y float32) {
	if math.IsNaN(float64(x)) {
		x = 0
	}
	if math.IsNaN(float64(y)) {
		y = 0
	}
	r.X = x
	r.Y = y
}

func (r *Rect) ModifyCenter(x, y float32) {
	if math.IsNaN(float64(x)) {
		x = 0
	}
	if math.IsNaN(float64(y)) {
		y = 0
	}
	r.X += x
	r.Y += y
}

func (r *Rect) Center() mgl32.Vec3 {
	return mgl32.Vec3{r.X, r.Y, 0}
}
