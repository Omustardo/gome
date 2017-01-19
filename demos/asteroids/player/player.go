package player

import (
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/demos/asteroids/bullet"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/util"
)

type Player struct {
	model.Model
	MoveSpeed     float32
	RotationSpeed mgl32.Quat

	// canFireAt is the next time the player is able to fire a weapon.
	canFireAt time.Time
	// FireRate is the duration between attacks.
	FireRate time.Duration
}

func New(mesh mesh.Mesh) *Player {
	return &Player{
		Model: model.Model{
			Mesh: mesh,
			Entity: entity.Entity{
				Position: mgl32.Vec3{0, -400, 0},
				// Rotate the model so it starts facing directly toward the positive Y axis, which is up on the user's screen.
				// Remember that these rotations are applied in the order specified by the final parameter
				// and that like a unit circle, positive values go to the "left" and negative values go to the "right".
				// X,Y,Z correspond to Roll, Pitch, and Yaw.
				Rotation: mgl32.AnglesToQuat(mgl32.DegToRad(90), mgl32.DegToRad(-90), 0, mgl32.XYZ),
				Scale:    mgl32.Vec3{50, 50, 50},
			},
		},
		MoveSpeed:     500,
		RotationSpeed: mgl32.AnglesToQuat(0, 0, mgl32.DegToRad(360/3), mgl32.XYZ),
		canFireAt:     time.Now(),
		FireRate:      time.Second * 1,
	}
}

func (p *Player) FireWeapon() *bullet.Bullet {
	if !p.CanFire() {
		return nil
	}
	p.canFireAt = p.canFireAt.Add(p.FireRate)
	b := bullet.New()

	// TODO: This is an odd forward vector - it's dependent on how the model is originally defined and loaded facing -X
	forward := p.Rotation.Rotate(mgl32.Vec3{-1, 0, 0}) // Forward is rotated to point in the direction the ship faces.
	b.Velocity = forward.Mul(800)
	b.Position = p.Position.Add(forward.Mul(p.Scale.X() * 1.3))
	return b
}

func (p *Player) CanFire() bool {
	return time.Now().After(p.canFireAt)
}

func (p *Player) Move(forward, back bool, delta float32) {
	var move float32
	if forward {
		move += delta
	}
	if back {
		move -= delta
	}
	moveSpeed := float32(500)

	// TODO: This is an odd forward vector - it's dependent on how the model is originally defined and loaded facing -X
	forwardDir := p.Rotation.Rotate(mgl32.Vec3{-1, 0, 0}) // Forward is rotated to point in the direction the ship faces.
	forwardDir = forwardDir.Mul(move * moveSpeed)         // direction * speed = distance
	p.ModifyPositionV(forwardDir)                         // current position + distance vector = final location
}

func (p *Player) Rotate(left, right bool, delta float32) {
	// rotate is the direction and amount to rotate.
	// Note that, like a unit circle, radians of higher positive value are toward the "left", while negative are to the "right".
	var rotationScale float32
	if left {
		rotationScale += delta
	}
	if right {
		rotationScale -= delta
	}
	if rotationScale != 0 {
		p.ModifyRotationGlobalQ(util.ScaleQuatRotation(p.RotationSpeed, rotationScale))
	}
}
