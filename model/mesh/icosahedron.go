package mesh

import (
	"math"

	"encoding/binary"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
)

func initializeIcosahedron() Mesh {
	faces := icosahedronFaces()

	vertices := make([]mgl32.Vec3, 0, 60)
	for _, f := range faces {
		vertices = append(vertices, f[0], f[1], f[2])
	}
	vertexBytes := bytecoder.Vec3(binary.LittleEndian, vertices...)
	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, vertexBytes, gl.STATIC_DRAW)

	normals := icosahedronNormals(faces)
	normalBytes := bytecoder.Vec3(binary.LittleEndian, normals...)
	normalVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, normalVBO)
	gl.BufferData(gl.ARRAY_BUFFER, normalBytes, gl.STATIC_DRAW)

	// TODO: add texture coordinates
	//texCoords := circleTexCoords(numCircleSegments)
	//texCoordBytes := bytecoder.Vec2(binary.LittleEndian, texCoords...)
	//texCoordsVBO := gl.CreateBuffer()
	//gl.BindBuffer(gl.ARRAY_BUFFER, texCoordsVBO)
	//gl.BufferData(gl.ARRAY_BUFFER, texCoordBytes, gl.STATIC_DRAW)

	return NewMesh(vertexVBO, gl.Buffer{}, normalVBO, gl.TRIANGLES, 20*3, nil, gl.Texture{}, gl.Buffer{})
}

// NewIcosahedron returns a mesh for a 20 sided figure.
// TODO: Implement texture coordinates for this figure. For now only color is supported.
func NewIcosahedron(col *color.NRGBA, texture gl.Texture) Mesh {
	if texture.Valid() {
		panic("texture mapping unsupported for NewIcosahedron")
	}
	ic := icosahedron
	ic.Color = col
	ic.SetTexture(texture)
	return ic
}

// icosahedronVertices returns a slice of the points making up the 20 faces of an icosahedron.
// The return value is 60 vertices, grouped into 20 triangles.
// Based on http://blog.andreaskahler.com/2009/06/creating-icosphere-mesh-in-code.html
func icosahedronFaces() [20][3]mgl32.Vec3 {
	t := float32((1.0 + math.Sqrt(5.0)) / 2.0)
	points := make([]mgl32.Vec3, 0, 12)
	points = append(points, mgl32.Vec3{-1, t, 0})
	points = append(points, mgl32.Vec3{1, t, 0})
	points = append(points, mgl32.Vec3{-1, -t, 0})
	points = append(points, mgl32.Vec3{1, -t, 0})
	points = append(points, mgl32.Vec3{0, -1, t})
	points = append(points, mgl32.Vec3{0, 1, t})
	points = append(points, mgl32.Vec3{0, -1, -t})
	points = append(points, mgl32.Vec3{0, 1, -t})
	points = append(points, mgl32.Vec3{t, 0, -1})
	points = append(points, mgl32.Vec3{t, 0, 1})
	points = append(points, mgl32.Vec3{-t, 0, -1})
	points = append(points, mgl32.Vec3{-t, 0, 1})
	for i := range points {
		points[i] = points[i].Normalize()
	}

	faces := [20][3]mgl32.Vec3{
		// 5 faces around point 0
		{points[0], points[11], points[5]},
		{points[0], points[5], points[1]},
		{points[0], points[1], points[7]},
		{points[0], points[7], points[10]},
		{points[0], points[10], points[11]},

		// 5 adjacent faces
		{points[1], points[5], points[9]},
		{points[5], points[11], points[4]},
		{points[11], points[10], points[2]},
		{points[10], points[7], points[6]},
		{points[7], points[1], points[8]},

		// 5 faces around point 3
		{points[3], points[9], points[4]},
		{points[3], points[4], points[2]},
		{points[3], points[2], points[6]},
		{points[3], points[6], points[8]},
		{points[3], points[8], points[9]},

		// 5 adjacent faces
		{points[4], points[9], points[5]},
		{points[2], points[4], points[11]},
		{points[6], points[2], points[10]},
		{points[8], points[6], points[7]},
		{points[9], points[8], points[1]},
	}
	return faces
}

// icosahedronNormals takes the triangles making up an icosahedron and returns a slice of normals.
// There will be one normal per vertex, for a total of 60 normals. All normals for a single triangle
// are the same.
func icosahedronNormals(faces [20][3]mgl32.Vec3) []mgl32.Vec3 {
	normals := make([]mgl32.Vec3, 0, 60)
	for _, f := range faces {
		// Cross product of two sides of a triangle is a surface normal.
		v := f[1].Sub(f[0])
		w := f[2].Sub(f[0])
		n := v.Cross(w)
		// Append the same normal three times, since we need 1 normal per vertex.
		normals = append(normals, n, n, n)
	}
	return normals
}
