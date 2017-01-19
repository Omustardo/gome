package asteroid

import (
	"math/rand"

	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/util"
)

const (
	_ = float32(100) * iota
	Small
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
	vel := mgl32.Vec3{rand.Float32(), rand.Float32(), 0}.Normalize().Mul(rand.Float32() * 150)
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
	a.ModifyPositionV(a.Velocity.Mul(deltaSeconds))

	// TODO: Wrap around position if it's too far from center/player
}

func (a *Asteroid) Split() (*Asteroid, *Asteroid) {
	if a == nil {
		log.Println("Attempted to split a nil asteroid")
		return nil, nil
	}
	// Asteroids start as copies of the original.
	a1, a2 := *a, *a
	// but move faster
	a1.Velocity = a.Velocity.Mul(1.3)
	a2.Velocity = a.Velocity.Mul(1.3)
	// and in slightly different directions
	a1.Velocity = mgl32.AnglesToQuat(0, 0, mgl32.DegToRad(30), mgl32.XYZ).Rotate(a1.Velocity)
	a2.Velocity = mgl32.AnglesToQuat(0, 0, mgl32.DegToRad(-30), mgl32.XYZ).Rotate(a2.Velocity)
	// and have slightly different rotations, just so they don't look identical
	a1.ModifyRotationLocal(mgl32.Vec3{0, 0, mgl32.DegToRad(30)})
	a2.ModifyRotationLocal(mgl32.Vec3{0, 0, mgl32.DegToRad(-30)})
	// and are moved slightly apart
	a1.ModifyPositionV(a1.Velocity.Normalize().Mul(a.Scale.X() * 0.6))
	a2.ModifyPositionV(a2.Velocity.Normalize().Mul(a.Scale.X() * 0.6))

	switch a.Scale.X() {
	case Large:
		a1.Scale = mgl32.Vec3{Medium, Medium, Medium}
		a2.Scale = mgl32.Vec3{Medium, Medium, Medium}
	case Medium:
		a1.Scale = mgl32.Vec3{Small, Small, Small}
		a2.Scale = mgl32.Vec3{Small, Small, Small}
	case Small:
		return nil, nil
	default:
		panic("unknown asteroid size")
	}
	return &a1, &a2
}
