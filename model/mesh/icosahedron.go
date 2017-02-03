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

	normals := triangleNormals(faces)
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
func icosahedronFaces() [][3]mgl32.Vec3 {
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

	faces := [][3]mgl32.Vec3{
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

// NewSubdividedIcosahedron returns a mesh for an n sided figure created by recursively dividing each face of an
// icosahedron (20 sided figure). Note that this recursive division of each triangular face grows very quickly and
// may cause unexpected crashes due to running out of memory if the number of divisions is too large.
func NewSubdividedIcosahedron(divisions int, col *color.NRGBA, texture gl.Texture) Mesh {
	if divisions < 0 {
		return Mesh{}
	}
	if divisions < len(subdividedIcosahedron) {
		if m, ok := subdividedIcosahedron[divisions]; ok {
			m.Color = col
			m.SetTexture(texture)
			return m
		}
	}

	if subdividedIcosahedron == nil {
		subdividedIcosahedron = make(map[int]Mesh)
	}
	faces := icosahedronFaces()

	// Make more spheres of increasing detail by subdividing the icosahedron faces,
	for i := 0; i <= divisions; i++ {
		// Only create the mesh if it hasn't been done yet.
		if _, ok := subdividedIcosahedron[i]; !ok {
			vertices := make([]mgl32.Vec3, 0, 20*3*int(math.Pow(float64(i), 4)))
			for _, f := range faces {
				vertices = append(vertices, f[0], f[1], f[2])
			}
			vertexBytes := bytecoder.Vec3(binary.LittleEndian, vertices...)
			vertexVBO := gl.CreateBuffer()

			gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
			gl.BufferData(gl.ARRAY_BUFFER, vertexBytes, gl.STATIC_DRAW)

			normals := triangleNormals(faces)
			normalBytes := bytecoder.Vec3(binary.LittleEndian, normals...)
			normalVBO := gl.CreateBuffer()
			gl.BindBuffer(gl.ARRAY_BUFFER, normalVBO)
			gl.BufferData(gl.ARRAY_BUFFER, normalBytes, gl.STATIC_DRAW)

			//texCoords := circleTexCoords(numCircleSegments)
			//texCoordBytes := bytecoder.Vec2(binary.LittleEndian, texCoords...)
			//texCoordsVBO := gl.CreateBuffer()
			//gl.BindBuffer(gl.ARRAY_BUFFER, texCoordsVBO)
			//gl.BufferData(gl.ARRAY_BUFFER, texCoordBytes, gl.STATIC_DRAW)

			subdividedIcosahedron[i] = NewMesh(vertexVBO, gl.Buffer{}, normalVBO, gl.TRIANGLES, len(vertices), nil, gl.Texture{}, gl.Buffer{})
		}
		// Divide each face into four faces and continue.
		newFaces := make([][3]mgl32.Vec3, 0, 4*len(faces))
		for _, f := range faces {
			// for each face, divide it into four triangles
			for _, face := range subdivideTriangle(f) {
				// and then push them outward so they are on the surface of the sphere
				for vert := range face {
					face[vert] = face[vert].Normalize()
				}
				newFaces = append(newFaces, face)
			}
		}
		faces = newFaces
	}

	if m, ok := subdividedIcosahedron[divisions]; ok {
		m.Color = col
		m.SetTexture(texture)
		return m
	}
	return Mesh{} // should never happen
}
