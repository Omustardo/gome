package util

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/model/line"
)

const axisLength = 1e12

var xAxis, yAxis, zAxis *line.Line

func init() {
	xAxis = &line.Line{
		P1:    mgl32.Vec3{-axisLength, 0, 0},
		P2:    mgl32.Vec3{axisLength, 0, 0},
		Color: &color.NRGBA{255, 0, 0, 255},
	}
	yAxis = &line.Line{
		P1:    mgl32.Vec3{0, -axisLength, 0},
		P2:    mgl32.Vec3{0, axisLength, 0},
		Color: &color.NRGBA{0, 255, 0, 255},
	}
	zAxis = &line.Line{
		P1:    mgl32.Vec3{0, 0, -axisLength},
		P2:    mgl32.Vec3{0, 0, axisLength},
		Color: &color.NRGBA{0, 0, 255, 255},
	}
}

// DrawXYZAxes is a utility function that draws the three basic X,Y,Z axes colored red, green, and blue respectively.
func DrawXYZAxes() {
	xAxis.Draw()
	yAxis.Draw()
	zAxis.Draw()
}
