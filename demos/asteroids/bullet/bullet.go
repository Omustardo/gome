package bullet

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
)

type Bullet struct {
	model.Model
	Velocity        mgl32.Vec3
	LifespanSeconds float32
}

func New() *Bullet {
	return &Bullet{
		Model: model.Model{
			Mesh: mesh.NewCircle(&color.NRGBA{200, 15, 15, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Position: mgl32.Vec3{},
				Scale:    mgl32.Vec3{15, 15, 0},
				Rotation: mgl32.QuatIdent(),
			},
		},
		Velocity:        mgl32.Vec3{},
		LifespanSeconds: 5.0,
	}
}

func (b *Bullet) Update(deltaSeconds float32) {
	b.ModifyPositionV(b.Velocity.Mul(deltaSeconds))
	b.LifespanSeconds -= deltaSeconds
}
