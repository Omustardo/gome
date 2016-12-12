package model

import (
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/shader"
)

type Mesh struct {
	VertexVBO, NormalVBO gl.Buffer
	TriangleCount        int
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

	gl.DrawArrays(gl.TRIANGLES, 0, 3*m.TriangleCount)
}
