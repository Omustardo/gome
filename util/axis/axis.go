package axis

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
)

const axisLength = 1e12

var xAxis, yAxis, zAxis *model.Model

// Initialize depends on mesh.Initialize being done first. TODO: This is the start of dependency hell. Stop it. Where could this package go to make it more logical? Inside of model?
func Initialize() {
	// Ignore the returned destroy functions on the line meshes since the axes should exist until the program ends.
	xMesh, _ := mesh.NewLine(
		mgl32.Vec3{-axisLength, 0, 0},
		mgl32.Vec3{axisLength, 0, 0},
		&color.NRGBA{255, 0, 0, 255},
	)
	yMesh, _ := mesh.NewLine(
		mgl32.Vec3{0, -axisLength, 0},
		mgl32.Vec3{0, axisLength, 0},
		&color.NRGBA{0, 255, 0, 255},
	)
	zMesh, _ := mesh.NewLine(
		mgl32.Vec3{0, 0, -axisLength},
		mgl32.Vec3{0, 0, axisLength},
		&color.NRGBA{0, 0, 255, 255},
	)

	xAxis = &model.Model{
		Mesh:   xMesh,
		Entity: entity.Default,
	}
	yAxis = &model.Model{
		Mesh:   yMesh,
		Entity: entity.Default,
	}
	zAxis = &model.Model{
		Mesh:   zMesh,
		Entity: entity.Default,
	}
}

// DrawXYZAxes is a utility function that draws the three basic X,Y,Z axes colored red, green, and blue respectively.
func DrawXYZAxes() {
	xAxis.Render()
	yAxis.Render()
	zAxis.Render()
}
