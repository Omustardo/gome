package shape

import (
	"encoding/binary"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/shader"
	"golang.org/x/mobile/exp/f32"
)

var _ Shape = (*Triangle)(nil)

type Triangle struct {
	P1, P2, P3 mgl32.Vec3
	R, G, B, A float32
}

func (t *Triangle) SetCenter(x, y float32) {
	center := t.Center()
	t.P1[0] += x - center[0]
	t.P2[0] += x - center[0]
	t.P3[0] += x - center[0]
	t.P1[1] += y - center[1]
	t.P2[1] += y - center[1]
	t.P3[1] += y - center[1]
}
func (t *Triangle) ModifyCenter(x, y float32) {
	t.P1[0] += x
	t.P2[0] += x
	t.P3[0] += x
	t.P1[1] += y
	t.P2[1] += y
	t.P3[1] += y
}

func (t *Triangle) Center() mgl32.Vec3 {
	x := (t.P1[0] + t.P2[0] + t.P3[0]) / 3
	y := (t.P1[1] + t.P2[1] + t.P3[1]) / 3
	return mgl32.Vec3{x, y, 0}
}

func (t *Triangle) Draw() {
	shader.Basic.SetDefaults()
	// TODO
	t.DrawFilled()
}

func (t *Triangle) DrawFilled() {
	shader.Basic.SetDefaults()
	shader.Basic.SetColor(t.R, t.G, t.B, t.A)

	// NOTE: Be careful of using len(vertices). It's NOT an array of floats - it's an array of bytes.
	vertices := f32.Bytes(binary.LittleEndian,
		t.P1.X(), t.P1.Y(), t.P1.Z(),
		t.P2.X(), t.P2.Y(), t.P2.Z(),
		t.P3.X(), t.P3.Y(), t.P3.Z(),
	)

	vbuffer := gl.CreateBuffer()                             // Generate buffer and returns a reference to it. https://www.khronos.org/opengles/sdk/docs/man/xhtml/glGenBuffers.xml
	gl.BindBuffer(gl.ARRAY_BUFFER, vbuffer)                  // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW) // store values in buffer

	itemSize := 3                                                 // because the points consist of 3 floats
	itemCount := 3                                                // number of vertices in total
	gl.EnableVertexAttribArray(shader.Basic.VertexPositionAttrib) // https://www.opengl.org/sdk/docs/man2/xhtml/glEnableVertexAttribArray.xml
	gl.VertexAttribPointer(shader.Basic.VertexPositionAttrib, itemSize, gl.FLOAT, false, 0, 0)

	gl.DrawArrays(gl.TRIANGLES, 0, itemCount)

	gl.DisableVertexAttribArray(shader.Basic.VertexPositionAttrib)
	gl.DeleteBuffer(vbuffer)
}
