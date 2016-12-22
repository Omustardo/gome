package asteroid

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/util"
)

const (
	Small = float32(8) * iota
	Medium
	Large
)

var Mesh mesh.Mesh

func SetMesh(m mesh.Mesh) {
	Mesh = m
}

type Asteroid struct {
	model.Model
	Velocity      mgl32.Vec3
	RotationSpeed mgl32.Vec3
}

func New() *Asteroid {
	pos := util.RandVec3().Normalize().Mul(1000) // 1000 == stage size. Need to pass this in or something.
	vel := util.RandVec3().Normalize()
	scale := Large
	rot := util.RandVec3().Normalize().Mul(2 * math.Pi)

	m := model.Model{
		Mesh: Mesh,
		Entity: entity.Entity{
			Position: pos,
			Scale:    mgl32.Vec3{scale, scale, scale},
			Rotation: rot,
		},
	}

	return &Asteroid{
		Model:         m,
		Velocity:      vel,
		RotationSpeed: util.RandVec3().Normalize().Mul(2 * math.Pi),
	}
}

func (a *Asteroid) Update(deltaSeconds float32) {
	a.Rotation.Dot(a.RotationSpeed.Mul(deltaSeconds))
	a.Position.Add(a.Velocity.Mul(deltaSeconds))

	// TODO: Wrap around position if it's too far from center/player
}
