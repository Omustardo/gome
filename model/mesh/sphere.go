package mesh

import (
	"encoding/binary"
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
)

// maxDetail determines how many times to recursively subdivide an icosahedron
func initializeSpheres(maxDetail int) {
	if spheres == nil {
		spheres = make(map[int]Mesh)
	}
	faces := icosahedronFaces()

	// Make more spheres of increasing detail by subdividing the icosahedron faces,
	for i := 0; i <= maxDetail; i++ {
		// Only create the mesh if it hasn't been done yet.
		if _, ok := spheres[i]; !ok {
			vertices := make([]mgl32.Vec3, 0, 20*3*int(math.Pow(float64(i), 4)))
			for _, f := range faces {
				vertices = append(vertices, f[0], f[1], f[2])
			}
			vertexBytes := bytecoder.Vec3(binary.LittleEndian, vertices...)
			vertexVBO := gl.CreateBuffer()

			gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
			gl.BufferData(gl.ARRAY_BUFFER, vertexBytes, gl.STATIC_DRAW)

			//texCoords := circleTexCoords(numCircleSegments)
			//texCoordBytes := bytecoder.Vec2(binary.LittleEndian, texCoords...)
			//texCoordsVBO := gl.CreateBuffer()
			//gl.BindBuffer(gl.ARRAY_BUFFER, texCoordsVBO)
			//gl.BufferData(gl.ARRAY_BUFFER, texCoordBytes, gl.STATIC_DRAW)

			// Use vertexVBO as the normalVBO to smooth out polygon edges.
			spheres[i] = NewMesh(vertexVBO, gl.Buffer{}, vertexVBO, gl.TRIANGLES, len(vertices), nil, gl.Texture{}, gl.Buffer{})
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
}

// NewSphere returns a sphere mesh. The higher the detail, the more triangles that make up the mesh.
// Making this number too large will very quickly make you run out of memory.
// Detail 0 is an icosahedron with 20 triangular faces. Each detail level has (detail ^ 4) * 3 * 20 vertices,
// so by detail 5 it's up to 37500 faces. Recommended value for a decently smooth sphere is 4.
func NewSphere(detail int, col *color.NRGBA, texture gl.Texture) Mesh {
	if _, ok := spheres[detail]; !ok {
		initializeSpheres(detail)
	}
	s := spheres[detail]
	s.Color = col
	s.SetTexture(texture)
	return s
}
