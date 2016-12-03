package shape

import (
	"encoding/binary"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/bytecoder"
)

var _ Shape = (*Line)(nil)

var (
	lineBuffer gl.Buffer
)

func loadLines() {
	lineBuffer = gl.CreateBuffer()
}

type Line struct {
	P1, P2     mgl32.Vec3
	R, G, B, A float32
}

func (l *Line) SetCenter(x, y float32) {
	center := l.Center()
	l.P1[0] += x - center[0]
	l.P2[0] += x - center[0]
	l.P1[1] += y - center[1]
	l.P2[1] += y - center[1]
}
func (l *Line) ModifyCenter(x, y float32) {
	l.P1[0] += x
	l.P2[0] += x
	l.P1[1] += y
	l.P2[1] += y
}

func (l *Line) Center() mgl32.Vec3 {
	x := (l.P1[0] + l.P2[0]) / 2
	y := (l.P1[1] + l.P2[1]) / 2
	return mgl32.Vec3{x, y, 0}
}

// Draw draws a line. It creates and deletes a buffer on the GPU to do this, whcih is relatively expensive.
// It's fine for drawing a few lines, but for many lines use a batched call.
func (l *Line) Draw() {
	shader.Basic.SetDefaults()
	gl.BindBuffer(gl.ARRAY_BUFFER, lineBuffer)
	vertices := bytecoder.Float32(binary.LittleEndian,
		l.P1[0], l.P1[1], l.P1[2],
		l.P2[0], l.P2[1], l.P2[2],
	)
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	itemSize := 3                                                 // we use vertices made up of 3 floats
	itemCount := 2                                                // 2 points
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	shader.Basic.SetColor(l.R, l.G, l.B, l.A)
	gl.DrawArrays(gl.LINES, 0, itemCount)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
}

// DrawFilled for a line is equivalent to Draw, but still required for the Shape interface.
func (l *Line) DrawFilled() {
	l.Draw()
}
