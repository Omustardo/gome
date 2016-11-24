package entity

import "github.com/go-gl/mathgl/mgl32"

type Entity interface {
	Center() mgl32.Vec3
}
