package mesh

import (
	"fmt"
	"image/color"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/glutil"
)

var (
	rect, wireframeRect     Mesh
	circle, wireframeCircle Mesh
	cube                    Mesh
	icosahedron             Mesh

	// subdividedIcosahedron is a mapping from detail level to a corresponding mesh.
	subdividedIcosahedron map[int]Mesh

	// spheres is a mapping from detail level to a corresponding mesh.
	// Initial detail levels are generated on Initialize() in initializeSpheres(), and higher detail levels are
	// created as needed in NewSphere().
	spheres map[int]Mesh // TODO: consider making this a slice.

	// Axes are the X,Y,Z axes in a mesh. All axes are of length 1, so they should be scaled before being drawn in most cases.
	// The XYZ axes are colored RGB respectively.
	Axes Mesh
)

// Loads models into buffers on the GPU. glfw.Init() must be called before calling this.
func Initialize() {
	fmt.Println("Generating Meshes...")
	initializeEmptyTexture()
	initializeEmptyBuffer()

	rect = initializeRect()
	cube = initializeCube()
	circle = initializeCircle()
	wireframeCircle = initializeWireframeCircle()
	wireframeRect = initializeWireframeRect()
	icosahedron = initializeIcosahedron()
	initializeSpheres(4)
	Axes = initializeAxes()
}

type Mesh struct {
	// Buffers are all private in order to set valid defaults for provided buffers that haven't been initialized.
	// For example, if the provided texture buffer is empty (no values provided), the setters replaces it with a default
	// texture containing all 0 values.
	// The purpose of these defaults is so the same shader can be used regardless of a few missing fields.

	// References to buffers on the GPU.
	vertices      gl.Buffer
	vertexIndices gl.Buffer

	normals gl.Buffer

	texture       gl.Texture
	textureCoords gl.Buffer

	// VBOMode is the gl Mode passed to a Draw call.
	// Most commonly, it is gl.TRIANGLES. See https://en.wikibooks.org/wiki/OpenGL_Programming/GLStart/Tut3
	vboMode gl.Enum

	// ItemCount is the number of items to be drawn.
	// For a rectangle to be drawn with gl.DrawArrays(gl.Triangles,...) this would be 2.
	// For a rectangle where only the edges are drawn with gl.LINE_LOOP, this would be 4.
	itemCount int

	// Color is 32-bit non-premultiplied RGBA. It is optional, but leaving it nil is the same as setting it to pure white.
	// Note that the color.Color interface's Color() function returns weird values (between 0 and 0xFFFF for avoiding overflow).
	// I recommend just accessing the RGBA fields directly.
	Color *color.NRGBA

	// BaseRotation is a rotation applied to the mesh.
	// This is intended to be used to orient the mesh so Up toward {0,1,0}, and Forward is towards {1,0,0}, since not all
	// meshes will be created in this orientation. Note that if you don't want to modify the default rotation, this must
	// be set to mgl32.QuatIdent().
	// Encouraging a consistent orientation for all meshes makes dealing with them much easier, but it is not absolutely
	// required.
	BaseRotation mgl32.Quat

	// TODO: Add a center value which is a added to position when rendering. As it is, position can be thought of as the
	// bottom left corner of a cube that bounds a mesh. Being able to change positioning to an arbitrary center point will be necessary.
}

// NewMeshFromArrays copies the input vertices, normals, and texture coordinates into buffers on the GPU.
func NewMeshFromArrays(vertices, normals []mgl32.Vec3, textureCoords []mgl32.Vec2) (Mesh, error) {
	var vertexBuffer, uvBuffer, normalBuffer gl.Buffer

	vertexBuffer = glutil.LoadBufferVec3(vertices)
	normalBuffer = glutil.LoadBufferVec3(normals)
	uvBuffer = glutil.LoadBufferVec2(textureCoords)

	if glError := gl.GetError(); glError != 0 {
		return Mesh{}, fmt.Errorf("gl.GetError: %v", glError)
	}

	return NewMesh(vertexBuffer, gl.Buffer{}, normalBuffer, gl.TRIANGLES, len(vertices), nil, gl.Texture{}, uvBuffer), nil
}

// NewMesh combines the input buffers and rendering information into a Mesh struct.
// Using this method requires loading OpenGL buffers yourself. It's not recommended for general use.
// Most standard use of meshes can be done via the standard ones (i.e. NewCube(), NewSphere(), NewRect())
// or by loading an model from file via the `asset` package.
func NewMesh(vertices, vertexIndices, normals gl.Buffer, vboMode gl.Enum, itemCount int, color *color.NRGBA, texture gl.Texture, textureCoords gl.Buffer) Mesh {
	if !vertices.Valid() {
		log.Println("Creating mesh with invalid vertex buffer")
	}
	m := Mesh{
		vertices:      vertices,
		vertexIndices: vertexIndices,
		vboMode:       vboMode,
		itemCount:     itemCount,
		Color:         color,
		BaseRotation:  mgl32.QuatIdent(),
	}
	m.SetNormalVBO(normals)
	m.SetTexture(texture)
	m.SetTextureCoords(textureCoords)
	return m
}

func (m *Mesh) VertexVBO() gl.Buffer {
	return m.vertices
}
func (m *Mesh) VertexIndices() gl.Buffer {
	return m.vertexIndices
}
func (m *Mesh) NormalVBO() gl.Buffer {
	return m.normals
}
func (m *Mesh) SetNormalVBO(normals gl.Buffer) {
	m.normals = normals
	if !m.normals.Valid() {
		m.normals = emptyBuffer
	}
}

// SetVBOMode allows changing what mode a mesh is rendered in. It's recommended not to use this on the built in meshes
// as they are created with the VBO Mode they expect to be rendered with.
// The VBO modes include gl.TRIANGLES, gl.LINES, gl.LINE_LOOP, etc.
func (m *Mesh) SetVBOMode(mode gl.Enum) {
	m.vboMode = mode
}
func (m *Mesh) VBOMode() gl.Enum {
	return m.vboMode
}
func (m *Mesh) ItemCount() int {
	return m.itemCount
}
func (m *Mesh) Texture() gl.Texture {
	return m.texture
}
func (m *Mesh) SetTexture(texture gl.Texture) {
	m.texture = texture
	if !m.texture.Valid() {
		m.texture = emptyTexture
	}
}
func (m *Mesh) TextureCoords() gl.Buffer {
	return m.textureCoords
}
func (m *Mesh) SetTextureCoords(coords gl.Buffer) {
	m.textureCoords = coords
	if !m.textureCoords.Valid() {
		m.textureCoords = emptyBuffer
	}
}

const (
	emptyBufferSize = 1024 * 1024 * 3
)

var (
	// TODO: Using these "empty" default buffers is ugly. I'm not sure how to do better though. If this is the way to go
	// make sure to document the size limitations as any meshes with undefined features that are larger than the size
	// limits will likely be hard to detect and diagnose.

	// EmptyTexture is a texture buffer filled with just four bytes: [255, 255, 255, 255]
	// It is meant to be used as in meshes that don't contain another texture.
	emptyTexture gl.Texture

	// EmptyNormals is a buffer on the GPU containing many zeros. This is used as a default for meshes that don't have the
	// relevant values defined, but still need to use a buffer since the shader requires having some values.
	// When used as a normal buffer, it means the vertices have no normal, which doesn't make logical sense, but seems to work in practice.
	// When used as a texture coordinate buffer in combination with emptyTexture, the zero values are used as coordinates and all reference a single pixel in the texture.
	emptyBuffer gl.Buffer
)

func initializeEmptyTexture() {
	var err error
	emptyTexture, err = glutil.LoadTextureData(1, 1, []uint8{255, 255, 255, 255}) // TODO: This sort of initialization behavior should be done in an actual init, where a panic is reasonable.
	if err != nil {
		panic(err)
	}
}

func initializeEmptyBuffer() {
	data := make([]float32, emptyBufferSize) // large array of 0 values
	emptyBuffer = glutil.LoadBufferFloat32(data)
}

// subdivideTriangle takes a triangle as input and returns the four triangles created by subdividing it.
func subdivideTriangle(tri [3]mgl32.Vec3) [][3]mgl32.Vec3 {
	// Get the three midpoints
	a := tri[0].Add(tri[1]).Mul(0.5)
	b := tri[1].Add(tri[2]).Mul(0.5)
	c := tri[2].Add(tri[0]).Mul(0.5)

	return [][3]mgl32.Vec3{
		{tri[0], a, c},
		{tri[1], b, a},
		{tri[2], c, b},
		{a, b, c},
	}
}

// TriangleNormals takes triangles that presumably make up mesh and returns a slice of normals.
// There will be one normal per vertex, for a total of 3 * len(triangles).
// All normals for a single triangle are the same.
// Note that the direction of the normals is based on the right hand rule, so the order by which the vertices are defined can flip the direction of the normal.
// If you need it in the other direction, just negate each normal:
// for i := range normals {
//   normals[i] = normals[i].Mul(-1)
// }
func TriangleNormals(triangles [][3]mgl32.Vec3) []mgl32.Vec3 {
	normals := make([]mgl32.Vec3, 0, 60)
	for _, f := range triangles {
		// Cross product of two sides of a triangle is a surface normal.
		v := f[1].Sub(f[0])
		w := f[2].Sub(f[0])
		n := v.Cross(w)
		// Append the same normal three times, since we need 1 normal per vertex.
		normals = append(normals, n, n, n)
	}
	return normals
}
