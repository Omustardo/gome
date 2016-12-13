package model

import (
	"image/color"

	"log"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/shader"
)

type Mesh struct {
	// VBO's are references to buffers on the GPU.
	VertexVBO, NormalVBO gl.Buffer

	// VBOMode is the gl Mode passed to a Draw call.
	// Most commonly, it is gl.TRIANGLES. See https://en.wikibooks.org/wiki/OpenGL_Programming/GLStart/Tut3
	VBOMode gl.Enum

	// ItemCount is the number of items to be drawn.
	// For a rectangle to be drawn with gl.DrawArrays(gl.Triangles,...) this would be 2.
	// For a rectangle where only the edges are drawn with gl.LINE_LOOP, this would be 4.
	ItemCount int

	// Color is 32-bit non-premultiplied RGBA. It is optional.
	Color *color.NRGBA

	Texture gl.Texture

	//	vboType      *gl.Enum // like gl.UNSIGNED_SHORT
}

type Model struct {
	Mesh
	entity.Entity
}

func (m *Model) Render() {
	shader.Model.SetDefaults()
	shader.Model.SetTranslationMatrix(m.Position.X(), m.Position.Y(), m.Position.Z())
	shader.Model.SetRotationMatrix(m.Rotation.X(), m.Rotation.Y(), m.Rotation.Z())
	shader.Model.SetScaleMatrix(m.Scale.X(), m.Scale.Y(), m.Scale.Z())

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Mesh.VertexVBO)
	gl.EnableVertexAttribArray(shader.Model.VertexPositionAttrib)
	gl.VertexAttribPointer(shader.Model.VertexPositionAttrib, 3, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Mesh.NormalVBO)
	gl.EnableVertexAttribArray(shader.Model.NormalAttrib)
	gl.VertexAttribPointer(shader.Model.NormalAttrib, 3, gl.FLOAT, false, 0, 0)

	switch m.VBOMode {
	case gl.TRIANGLES:
		gl.DrawArrays(gl.TRIANGLES, 0, 3*m.ItemCount)
	case gl.LINE_LOOP:
		gl.DrawArrays(gl.LINE_LOOP, 0, m.ItemCount)
	default:
		log.Printf("uknown VBO Mode: %v", m.VBOMode)
	}
}
