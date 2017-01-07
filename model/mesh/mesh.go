package mesh

import (
	"encoding/binary"
	"image/color"

	"log"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

var (
	line                    Mesh
	rect, wireframeRect     Mesh
	circle, wireframeCircle Mesh
	cube                    Mesh
	// TODO: Sphere, wireframeCube
)

// Loads models into buffers on the GPU. glfw.Init() must be called before calling this.
func Initialize() {
	initializeEmptyTexture()
	initializeEmptyTextureCoords()
	initializeEmptyBuffer()

	line = initializeLine()
	rect = initializeRect()
	cube = initializeCube()
	circle = initializeCircle()
	wireframeCircle = initializeWireframeCircle()
	wireframeRect = initializeWireframeRect()
	// initializeTriangle()
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

	// VBOMode is the gl Mode passed to a Draw call.
	// Most commonly, it is gl.TRIANGLES. See https://en.wikibooks.org/wiki/OpenGL_Programming/GLStart/Tut3
	vboMode gl.Enum

	// ItemCount is the number of items to be drawn.
	// For a rectangle to be drawn with gl.DrawArrays(gl.Triangles,...) this would be 2.
	// For a rectangle where only the edges are drawn with gl.LINE_LOOP, this would be 4.
	itemCount int

	// Color is 32-bit non-premultiplied RGBA. It is optional, but leaving it unset is the same as setting it to (1,1,1,1).
	// Note that the color.Color interface's Color() function returns weird values (between 0 and 0xFFFF for avoiding overflow).
	// I recommend just accessing the RGBA fields directly.
	Color *color.NRGBA

	texture       gl.Texture
	textureCoords gl.Buffer
}

// NewMesh combines the input buffers and rendering information into a Mesh struct.
// Using this method requires loading OpenGL buffers yourself. It's not recommended for general use.
// Most standard use of meshes can be done via the standard ones (i.e. NewCube(), NewSphere(), NewRect())
// or by loading an object file via the `asset` package.
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
		m.normals = EmptyNormals
	}
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
		m.texture = EmptyTexture
	}
}
func (m *Mesh) TextureCoords() gl.Buffer {
	return m.textureCoords
}
func (m *Mesh) SetTextureCoords(coords gl.Buffer) {
	m.textureCoords = coords
	if !m.textureCoords.Valid() {
		m.textureCoords = EmptyTextureCoords
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
	// It is referenced by EmptyTextureCoords and is meant to be used as in meshes that don't contain another texture.
	EmptyTexture gl.Texture

	// EmptyTextureCoords is a buffer on the GPU containing many 1's. These are all texture coordinates, pointing to the
	// single texel in EmptyTexture. This allows it to be used as a stand in / default texture mapping for objects that
	// don't have a texture of their own.
	EmptyTextureCoords gl.Buffer

	// EmptyNormals is a buffer on the GPU containing many zeros. This is used as a default for meshes that don't have the
	// relevant values defined, but still need to use a buffer since the shader requires having some values.
	EmptyNormals gl.Buffer
)

func initializeEmptyTexture() {
	EmptyTexture = gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, EmptyTexture)
	// NOTE: gl.FLOAT isn't enabled for texture data types unless gl.getExtension('OES_texture_float'); is set, so just use gl.UNSIGNED_BYTE
	//   See http://stackoverflow.com/questions/23124597/storing-floats-in-a-texture-in-opengl-es  http://stackoverflow.com/questions/22666556/webgl-texture-creation-trouble
	gl.TexImage2D(gl.TEXTURE_2D, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, []uint8{255, 255, 255, 255})
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture{}) // bind to "null" to prevent using the wrong texture by mistake.
}

func initializeEmptyTextureCoords() {
	coords := make([]float32, emptyBufferSize) // large array of 1 values
	for i := range coords {
		coords[i] = 1.0 // @@@@@@@@@@ why use 1 rather than 0? Thinking it should be 0.
	}
	textureCoordinates := bytecoder.Float32(binary.LittleEndian, coords...) // @@@@@@@@@@@@@@@@ TODO: Can this buffer be removed and just use EmptyBuffer instead?

	EmptyTextureCoords = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, EmptyTextureCoords)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)
}

func initializeEmptyBuffer() {
	data := make([]float32, emptyBufferSize) // large array of 0 values
	EmptyNormals = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, EmptyNormals)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, data...), gl.STATIC_DRAW)
}

// DestroyFunc is used to clear the buffers used by a mesh. Generally their call should be deferred.
type DestroyFunc func()
