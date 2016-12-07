package util

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/core/drawable"
	"github.com/omustardo/gome/geom"
)

const axisLength = 1e12

var xAxis, yAxis, zAxis *geom.Line

func init() {
	xAxis = &geom.Line{
		P1: mgl32.Vec3{-axisLength, 0, 0},
		P2: mgl32.Vec3{axisLength, 0, 0},
		Drawable: drawable.Drawable{
			Color: &color.RGBA{255, 0, 0, 255},
		},
	}
	yAxis = &geom.Line{
		P1: mgl32.Vec3{0, -axisLength, 0},
		P2: mgl32.Vec3{0, axisLength, 0},
		Drawable: drawable.Drawable{
			Color: &color.RGBA{0, 255, 0, 255},
		},
	}
	zAxis = &geom.Line{
		P1: mgl32.Vec3{0, 0, -axisLength},
		P2: mgl32.Vec3{0, 0, axisLength},
		Drawable: drawable.Drawable{
			Color: &color.RGBA{0, 0, 255, 255},
		},
	}
}

// DrawXYZAxes is a utility function that simply draws the three basic X,Y,Z axes
func DrawXYZAxes() {
	xAxis.Draw()
	yAxis.Draw()
	zAxis.Draw()
}
