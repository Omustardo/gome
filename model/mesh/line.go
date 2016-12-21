package mesh

import (
	"encoding/binary"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

// glfw.Init() must be called before calling this.
func initializeLine() Mesh {
	return NewMesh(gl.Buffer{}, gl.Buffer{}, gl.Buffer{}, gl.LINES, 2, nil, gl.Texture{}, gl.Buffer{})
}

func NewLine(p1, p2 mgl32.Vec3, col *color.NRGBA) (Mesh, DestroyFunc) {
	l := line
	l.Color = col

	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian,
		p1.X(), p1.Y(), p1.Z(),
		p2.X(), p2.Y(), p2.Z(),
	), gl.STATIC_DRAW)
	l.vertexVBO = vertexVBO

	SetValidDefaults(&l)
	return l, func() { gl.DeleteBuffer(vertexVBO) }
}
