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
	Small = float32(1) * iota
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
	pos := util.RandVec3().Normalize().Mul(100) // TODO: pass in limits on starting location
	pos[2] = 0
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
		RotationSpeed: util.RandVec3().Normalize().Mul(math.Pi / 3),
	}
}

func (a *Asteroid) Update(deltaSeconds float32) {
	a.Rotation[0] = a.Rotation[0] + a.RotationSpeed[0]*deltaSeconds
	a.Rotation[1] = a.Rotation[1] + a.RotationSpeed[1]*deltaSeconds
	a.Rotation[2] = a.Rotation[2] + a.RotationSpeed[2]*deltaSeconds
	a.Position.Add(a.Velocity.Mul(deltaSeconds))

	// TODO: Wrap around position if it's too far from center/player
}
