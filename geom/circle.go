package geom

import (
	"encoding/binary"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/drawable"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/bytecoder"
)

const numCircleSegments = 360

var (
	// Buffers hold the float32 coordinates of vertices that make up a circle, converted to a byte array.
	// This is the format required by OpenGL vertex buffers. These two buffers are used for all circles by modifying
	// the Scale, Rotation, and Translation matrices in the shader.
	circleTrianglesBuffer gl.Buffer
	circleLineLoopBuffer  gl.Buffer
)

func initializeCircle() {
	// Generates triangles to make a full circle entered at [0,0]. Not just the edges.
	tmp := mgl32.Circle(1.0, 1.0, numCircleSegments)

	// The values are good as is for use with draw mode gl.TRIANGLES.
	circleTriangleSegmentVertices := bytecoder.Vec2(binary.LittleEndian, tmp...)
	circleTrianglesBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, circleTrianglesBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, circleTriangleSegmentVertices, gl.STATIC_DRAW)

	// To get the line segment vertices, just ignore the first of every trio since that's the center.
	lineSegments := make([]mgl32.Vec2, numCircleSegments*2)
	for i := 0; i < numCircleSegments; i++ {
		lineSegments[i*2], lineSegments[i*2+1] = tmp[i*3+1], tmp[i*3+2]
	}
	circleLineSegmentVertices := bytecoder.Vec2(binary.LittleEndian, lineSegments...)
	circleLineLoopBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, circleLineLoopBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, circleLineSegmentVertices, gl.STATIC_DRAW)
}

type Circle struct {
	entity.Entity
	drawable.Attributes
}

func (c *Circle) DrawWireframe() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(c.Color)
	shader.Basic.SetTranslationMatrix(c.Position.X(), c.Position.Y(), c.Position.Z())
	shader.Basic.SetScaleMatrix(c.Scale.X(), c.Scale.Y(), c.Scale.Z())

	gl.BindBuffer(gl.ARRAY_BUFFER, circleLineLoopBuffer)
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	itemSize := 2                                                 // we use vertices made up of 2 floats
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	itemCount := numCircleSegments * 2 // 2 vertices per segment
	gl.DrawArrays(gl.LINE_LOOP, 0, itemCount)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}

func (c *Circle) DrawFilled() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(c.Color)
	shader.Basic.SetTranslationMatrix(c.Position.X(), c.Position.Y(), 0)
	shader.Basic.SetScaleMatrix(c.Scale.X(), c.Scale.Y(), c.Scale.Z())

	gl.BindBuffer(gl.ARRAY_BUFFER, circleTrianglesBuffer)
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	itemSize := 2                                                 // we use vertices made up of 2 floats
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	itemCount := numCircleSegments // One triangle per segment
	gl.DrawArrays(gl.TRIANGLES, 0, itemCount)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}
