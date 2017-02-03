package mesh

import (
	"encoding/binary"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
)

// TODO: Line is an odd case. Other meshes can be easily scaled and rotated while still using the same basic built in
// model, but that doesn't work for lines. The way this is, each line requires creating a buffer on the GPU which is
// never deleted. Adding a delete method just for line meshes is odd, since none of the other models do it, but it's
// probably necessary.
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
