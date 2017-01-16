package asteroid

import (
	"math/rand"

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
	RotationSpeed mgl32.Quat
}

func New() *Asteroid {
	pos := util.RandVec3().Normalize().Mul(100) // TODO: pass in limits on starting location
	pos[2] = 0
	vel := mgl32.Vec3{rand.Float32(), rand.Float32(), 0}.Normalize().Mul(rand.Float32() * 15)
	scale := Large

	m := model.Model{
		Mesh: Mesh,
		Entity: entity.Entity{
			Position: pos,
			Scale:    mgl32.Vec3{scale, scale, scale},
			Rotation: util.RandQuat(),
		},
	}

	return &Asteroid{
		Model:         m,
		Velocity:      vel,
		RotationSpeed: util.ScaleQuatRotation(util.RandQuat(), 0.8), // Scale down rotation speed so it isn't too fast.
	}
}

func (a *Asteroid) Update(deltaSeconds float32) {
	a.ModifyRotationLocalQ(util.ScaleQuatRotation(a.RotationSpeed, deltaSeconds))
	a.ModifyCenterV(a.Velocity.Mul(deltaSeconds))

	// TODO: Wrap around position if it's too far from center/player
}
