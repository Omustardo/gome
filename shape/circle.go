package shape

import (
	"encoding/binary"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/bytecoder"
)

const numCircleSegments = 360

var _ Shape = (*Circle)(nil)

var (
	// Buffers hold the float32 coordinates of vertices that make up a circle, converted to a byte array.
	// This is the format required by OpenGL vertex buffers. These two buffers are used for all circles by modifying
	// the Scale, Rotation, and Translation matrices in the shader.
	circleTriangleSegmentBuffer gl.Buffer
	circleLineSegmentBuffer     gl.Buffer
)

func loadCircles() {
	// Generates triangles to make a full circle entered at [0,0]. Not just the edges.
	tmp := mgl32.Circle(1.0, 1.0, numCircleSegments)

	// The values are good as is for making triangles.
	circleTriangleSegmentVertices := bytecoder.Vec2(binary.LittleEndian, tmp...)
	circleTriangleSegmentBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, circleTriangleSegmentBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, circleTriangleSegmentVertices, gl.STATIC_DRAW)

	// To get the line segment vertices, just ignore the first of every trio since that's the center.
	lineSegments := make([]mgl32.Vec2, numCircleSegments*2)
	for i := 0; i < numCircleSegments; i++ {
		lineSegments[i*2], lineSegments[i*2+1] = tmp[i*3+1], tmp[i*3+2]
	}
	circleLineSegmentVertices := bytecoder.Vec2(binary.LittleEndian, lineSegments...)
	circleLineSegmentBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, circleLineSegmentBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, circleLineSegmentVertices, gl.STATIC_DRAW)
}

type Circle struct {
	Pos        mgl32.Vec3
	Radius     float32
	R, G, B, A float32
}

func (c *Circle) SetCenter(x, y float32) {
	if math.IsNaN(float64(x)) {
		x = 0
	}
	if math.IsNaN(float64(y)) {
		y = 0
	}
	c.Pos[0], c.Pos[1] = x, y
}

func (c *Circle) ModifyCenter(x, y float32) {
	if math.IsNaN(float64(x)) {
		x = 0
	}
	if math.IsNaN(float64(y)) {
		y = 0
	}
	c.Pos[0] += x
	c.Pos[1] += y
}

func (c *Circle) Center() mgl32.Vec3 {
	return c.Pos
}

func (c *Circle) Draw() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(c.R, c.G, c.B, c.A)
	shader.Basic.SetTranslationMatrix(c.Pos.X(), c.Pos.Y(), 0)
	shader.Basic.SetScaleMatrix(c.Radius, c.Radius, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, circleLineSegmentBuffer)
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	itemSize := 2                                                 // we use vertices made up of 2 floats
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	itemCount := numCircleSegments * 2 // 2 vertices per segment
	gl.DrawArrays(gl.LINE_LOOP, 0, itemCount)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}

func (c *Circle) DrawFilled() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(c.R, c.G, c.B, c.A)
	shader.Basic.SetTranslationMatrix(c.Pos.X(), c.Pos.Y(), 0)
	shader.Basic.SetScaleMatrix(c.Radius, c.Radius, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, circleTriangleSegmentBuffer)
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	itemSize := 2                                                 // we use vertices made up of 2 floats
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	itemCount := numCircleSegments // One triangle per segment
	gl.DrawArrays(gl.TRIANGLES, 0, itemCount)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}
