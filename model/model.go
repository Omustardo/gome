package model

import (
	"image/color"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
)

type Model struct {
	// Tag is a human readable string for debugging.
	// NOTE: If this field starts being used for anything but human use, something is wrong.
	Tag string

	// Hidden determines whether the model will be rendered or not. False by default.
	Hidden bool

	mesh.Mesh
	entity.Entity
}

func (m *Model) Render() {
	if m == nil {
		log.Panic("Attempted to draw a nil model") // TODO: Not fatal error with better logging
		return
	}
	if m.Hidden || (m.Mesh == mesh.Mesh{}) {
		return
	}
	if m.Scale.X() == 0 && m.Scale.Y() == 0 && m.Scale.Z() == 0 {
		log.Println("Attempted to draw a model with scale [0,0,0]")
		return
	}
	if !m.Mesh.VertexVBO().Valid() {
		log.Println("Attempted to draw a model with no vertices")
		return
	}

	// TODO: Consider a "modelviewer" feature - let meshes have their own Render() method where they are rendered within a unit cube centered at the origin with no lighting or other world effects.
	shader.Model.SetDefaults()
	shader.Model.SetTranslationMatrix(m.Position.X(), m.Position.Y(), m.Position.Z())
	shader.Model.SetRotationMatrixQ(m.Rotation)
	shader.Model.SetScaleMatrix(m.Scale.X(), m.Scale.Y(), m.Scale.Z())
	shader.Model.SetColor(m.Mesh.Color)
	shader.Model.SetTexture(m.Mesh.Texture())

	gl.BindBuffer(gl.ARRAY_BUFFER, m.Mesh.VertexVBO())
	gl.EnableVertexAttribArray(shader.Model.VertexPositionAttrib) // TODO: Can these VertexAttribArrays be enabled a single time in shader initialization and then just always used?
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

func (m *Model) RenderRotationAxes() {
	// TODO
}

// RenderDebugSphere draws three circles that live in the bounding sphere of the model.
func (m *Model) RenderDebugSphere() {
	r := Model{
		Mesh:   mesh.NewCircleOutline(&color.NRGBA{255, 0, 0, 255}),
		Entity: m.Entity,
	}
	r.Render()
	r.ModifyRotationLocal(mgl32.Vec3{0, mgl32.DegToRad(90), 0})
	r.Color = &color.NRGBA{0, 255, 0, 255}
	r.Render()
	r.ModifyRotationLocal(mgl32.Vec3{mgl32.DegToRad(90), 0, 0})
	r.Color = &color.NRGBA{0, 0, 255, 255}
	r.Render()
}

// TODO: Global shader variables that keeps track of current programs/other vars so we don't have so many shader.UseProgram() and others.
// This sort of thing really ought to be in goxjs/gl rather than in here...

// TODO: Make sure there's documentation that all Index Buffers must be of type gl.UNSIGNED_SHORT (golang's uint16)
// and all other buffers must be gl.Float (golang's float32).
