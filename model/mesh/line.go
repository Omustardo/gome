package mesh

import (
	"encoding/binary"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

func NewLine(p1, p2 mgl32.Vec3, col *color.NRGBA) Mesh {
	vertexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian,
		p1.X(), p1.Y(), p1.Z(),
		p2.X(), p2.Y(), p2.Z(),
	), gl.STATIC_DRAW)

	line := NewMesh(vertexBuffer, gl.Buffer{}, gl.Buffer{}, gl.LINES, 2, nil, gl.Texture{}, gl.Buffer{})
	line.Color = col
	return line
}
