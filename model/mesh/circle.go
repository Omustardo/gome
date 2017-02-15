package mesh

import (
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/glutil"
)

const numCircleSegments = 360

func initializeCircle() Mesh {
	vertexVBO := glutil.LoadBufferVec3(circleVertices(numCircleSegments))
	texCoordsVBO := glutil.LoadBufferVec2(circleTexCoords(numCircleSegments))

	// item count is numSegments+2 because it's the total number of vertices in the fan:
	// one for the center, and one for each point on the circle, and then a single duplicate to close the circle.
	return NewMesh(vertexVBO, gl.Buffer{}, gl.Buffer{}, gl.TRIANGLE_FAN, numCircleSegments+2, nil, gl.Texture{}, texCoordsVBO)
}

func initializeWireframeCircle() Mesh {
	vertexVBO := glutil.LoadBufferVec3(wireframeCircleVertices(numCircleSegments))
	return NewMesh(vertexVBO, gl.Buffer{}, gl.Buffer{}, gl.LINE_LOOP, numCircleSegments, nil, gl.Texture{}, gl.Buffer{})
}

func NewCircle(col *color.NRGBA, texture gl.Texture) Mesh {
	c := circle
	c.Color = col
	c.SetTexture(texture)
	return c
}

func NewCircleOutline(col *color.NRGBA) Mesh {
	c := wireframeCircle
	c.Color = col
	return c
}

// Based on the github.com/go-gl/mathgl/mgl32 Circle function, but uses Vec3 and returns vertices suitable for a gl.TRIANGLE_FAN.
func circleVertices(numSlices int) []mgl32.Vec3 {
	radius := float32(1.0)
	twoPi := float32(2.0 * math.Pi)

	circlePoints := make([]mgl32.Vec3, 0, numSlices+2)
	// Add the center of the circle.
	circlePoints = append(circlePoints, mgl32.Vec3{0, 0, 0})

	step := twoPi / float32(numSlices)
	for i := 0; i < numSlices; i++ {
		currRadians := float32(i) * step
		sin, cos := math.Sincos(float64(currRadians))
		circlePoints = append(circlePoints, mgl32.Vec3{float32(cos) * radius, float32(sin) * radius, 0})
	}
	// Now add the final point to finish the last triangle.
	circlePoints = append(circlePoints, mgl32.Vec3{radius, 0, 0})
	return circlePoints
}

// Generates texture coordinates for the gl.TRIANGLE_FAN based vertices generated by circleVertices().
func circleTexCoords(numSlices int) []mgl32.Vec2 {
	// Since texture coordinates range from 0 to 1, we must consider the coordinates of a circle centered at (0.5, 0.5)
	// with radius 0.5.
	radius := float32(0.5)
	center := mgl32.Vec2{0.5, 0.5}

	twoPi := float32(2.0 * math.Pi)

	circlePoints := make([]mgl32.Vec2, 0, numSlices+2)
	// Add the center of the circle.
	circlePoints = append(circlePoints, center)

	step := twoPi / float32(numSlices)
	for i := 0; i < numSlices; i++ {
		currRadians := float32(i) * step
		sin, cos := math.Sincos(float64(currRadians))
		circlePoints = append(circlePoints, mgl32.Vec2{float32(cos) * radius, float32(sin) * radius}.Add(center))
	}
	// Now add the final point to finish the last triangle.
	circlePoints = append(circlePoints, mgl32.Vec2{radius, 0}.Add(center))
	return circlePoints
}

func wireframeCircleVertices(numSlices int) []mgl32.Vec3 {
	radius := float32(1.0)
	twoPi := float32(2.0 * math.Pi)
	circlePoints := make([]mgl32.Vec3, 0, numSlices+2)

	step := twoPi / float32(numSlices)
	for i := 0; i < numSlices; i++ {
		currRadians := float32(i) * step
		sin, cos := math.Sincos(float64(currRadians))
		circlePoints = append(circlePoints, mgl32.Vec3{float32(cos) * radius, float32(sin) * radius, 0})
	}
	return circlePoints
}
