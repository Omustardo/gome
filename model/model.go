package model

import (
	"image/color"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
)

type Model struct {
	Tag string
	mesh.Mesh
	entity.Entity
}

func (m *Model) Render() {

	// TODO: Consider not using this. It's an inexpensive call, but doing it for every model every frame is a bit much.
	// SetValidDefaults makes sure that all of the buffers are set, and sets them to the default "Empty" buffers if not.
	// The only reason this would be needed is if the functions to create meshes forgot to call it before returning, or
	// if users are directly messing with mesh internals, which they shouldn't be able to do with everything being private.
	mesh.SetValidDefaults(&m.Mesh)

	// TODO: Consider a "modelviewer" feature - let meshes have their own Render() method where they are rendered within a unit cube centered at the origin with no lighting or other world effects.
	shader.Model.SetDefaults()
	shader.Model.SetTranslationMatrix(m.Position.X(), m.Position.Y(), m.Position.Z())
	shader.Model.SetRotationMatrix(m.Rotation.X(), m.Rotation.Y(), m.Rotation.Z())
	shader.Model.SetScaleMatrix(m.Scale.X(), m.Scale.Y(), m.Scale.Z())
	shader.Model.SetColor(m.Mesh.Color)
	shader.Model.SetAmbientLight(&color.NRGBA{255, 255, 255, 0}) // &color.NRGBA{60, 60, 60, 0})
	shader.Model.SetTexture(m.Mesh.Texture())

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Mesh.VertexVBO())
	gl.EnableVertexAttribArray(shader.Model.VertexPositionAttrib)
	gl.VertexAttribPointer(shader.Model.VertexPositionAttrib, 3, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Mesh.NormalVBO())
	gl.EnableVertexAttribArray(shader.Model.NormalAttrib)
	gl.VertexAttribPointer(shader.Model.NormalAttrib, 3, gl.FLOAT, false, 0, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Mesh.TextureCoords())
	gl.EnableVertexAttribArray(shader.Model.TextureCoordAttrib)
	gl.VertexAttribPointer(shader.Model.TextureCoordAttrib, 2, gl.FLOAT, false, 0, 0)

	if m.VertexIndices().Valid() {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, m.Mesh.VertexIndices())
		gl.DrawElements(m.VBOMode(), m.ItemCount(), gl.UNSIGNED_SHORT, 0)
	} else {
		gl.DrawArrays(m.VBOMode(), 0, m.ItemCount())
	}
}

// TODO: Global shader variables that keeps track of current programs/other vars so we don't have so many shader.UseProgram() and others.
// This sort of thing really ought to be in goxjs/gl rather than in here...

// TODO: Make sure there's documentation that all Index Buffers must be of type gl.UNSIGNED_SHORT (golang's uint16)
// and all other buffers must be gl.Float (golang's float32).
