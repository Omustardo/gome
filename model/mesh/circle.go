package mesh

import (
	"encoding/binary"
	"image/color"

	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

const numCircleSegments = 36

func initializeCircle() Mesh {
	// Generates triangles to make a full circle entered at [0,0,0]. Not just the edges.
	tmp := circleVertices(1.0, 1.0, numCircleSegments)

	// The values are good as is for use with draw mode gl.TRIANGLES.
	vertices := bytecoder.Vec3(binary.LittleEndian, tmp...)
	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	//// To get the line segment vertices, just ignore the first of every trio since that's the center.
	//lineSegments := make([]mgl32.Vec2, numCircleSegments*2)
	//for i := 0; i < numCircleSegments; i++ {
	//	lineSegments[i*2], lineSegments[i*2+1] = tmp[i*3+1], tmp[i*3+2]
	//}
	//circleLineSegmentVertices := bytecoder.Vec2(binary.LittleEndian, lineSegments...)
	//circleLineLoopBuffer = gl.CreateBuffer()
	//gl.BindBuffer(gl.ARRAY_BUFFER, circleLineLoopBuffer)
	//gl.BufferData(gl.ARRAY_BUFFER, circleLineSegmentVertices, gl.STATIC_DRAW)

	// TODO: Need a default texture coordinate mapping for circles.
	// item count is numSegments+2 because it's the total number of vertices in the fan:
	// one for the center, and one for each point on the circle, and then a single duplicate to close the circle.
	return NewMesh(vertexVBO, gl.Buffer{}, gl.Buffer{}, gl.TRIANGLE_FAN, numCircleSegments+2, nil, gl.Texture{}, gl.Buffer{})
}

func NewCircle(color *color.NRGBA, texture gl.Texture) Mesh {
	c := circle
	c.Color = color
	c.texture = texture
	if !c.texture.Valid() {
		c.texture = EmptyTexture
	}
	return c
}

// Based on the github.com/go-gl/mathgl/mgl32 Circle function, but uses Vec3 and returns vertices suitable for a triangle fan.
func circleVertices(radiusX, radiusY float32, numSlices int) []mgl32.Vec3 {
	twoPi := float32(2.0 * math.Pi)

	circlePoints := make([]mgl32.Vec3, 0, numSlices+2)
	// Add the center of the circle.
	circlePoints = append(circlePoints, mgl32.Vec3{0, 0, 0})

	step := twoPi / float32(numSlices)
	for i := 0; i < numSlices; i++ {
		currRadians := float32(i) * step
		sin, cos := math.Sincos(float64(currRadians))
		circlePoints = append(circlePoints, mgl32.Vec3{float32(cos) * radiusX, float32(sin) * radiusY, 0})
	}
	// Now add the final point to finish the last triangle.
	circlePoints = append(circlePoints, mgl32.Vec3{radiusX, 0, 0})
	return circlePoints
}
