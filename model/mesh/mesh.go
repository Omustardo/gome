package mesh

import (
	"encoding/binary"
	"image/color"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/bytecoder"
)

var Line, rect, cube, circle, wireframeRect Mesh // TODO: Sphere. Consider removing WireframeRect

// Loads models into buffers on the GPU. glfw.Init() must be called before calling this.
func Initialize() {
	initializeEmptyTexture()
	initializeEmptyTextureCoords()
	initializeEmptyNormals()
	initializeIndices()

	rect = initializeRect()
	cube = initializeCube()
	circle = initializeCircle()
	wireframeRect = initializeWireframeRect()
	// initializeTriangle()
}

type Mesh struct {
	// usesElementArray determines whether rendering should use DrawElements, rather than the basic DrawArrays.
	usesElementArray bool

	// VBO's are references to buffers on the GPU.
	vertexVBO     gl.Buffer
	vertexIndices gl.Buffer

	normalVBO gl.Buffer

	// VBOMode is the gl Mode passed to a Draw call.
	// Most commonly, it is gl.TRIANGLES. See https://en.wikibooks.org/wiki/OpenGL_Programming/GLStart/Tut3
	vboMode gl.Enum

	// ItemCount is the number of items to be drawn.
	// For a rectangle to be drawn with gl.DrawArrays(gl.Triangles,...) this would be 2.
	// For a rectangle where only the edges are drawn with gl.LINE_LOOP, this would be 4.
	itemCount int

	// Color is 32-bit non-premultiplied RGBA. It is optional.
	Color *color.NRGBA

	texture       gl.Texture
	textureCoords gl.Buffer

	//	vboType      *gl.Enum // like gl.UNSIGNED_SHORT
	//
}

// NewMesh combines the input buffers and rendering information into a Mesh struct.
// Using this method requires loading OpenGL buffers yourself. It's not recommended for general use.
// Most standard use of meshes can be done via the standard ones (i.e. NewCube(), NewSphere(), NewRect())
// or by loading an object file via the `asset` package.
// TODO: Consider returning an error if vertex or normal VBO's are invalid (or any other invalid items are in the input)
func NewMesh(vertexVBO, vertexIndices, normalVBO gl.Buffer, vboMode gl.Enum, itemCount int, color *color.NRGBA, texture gl.Texture, textureCoords gl.Buffer) Mesh {
	if !texture.Valid() {
		texture = EmptyTexture
	}
	if !textureCoords.Valid() {
		textureCoords = EmptyTextureCoords
	}
	if !normalVBO.Valid() {
		normalVBO = EmptyNormals
	}
	return Mesh{
		vertexVBO:     vertexVBO,
		vertexIndices: vertexIndices,
		normalVBO:     normalVBO,
		vboMode:       vboMode,
		itemCount:     itemCount,
		Color:         color,
		texture:       texture,
		textureCoords: textureCoords,
	}
}

func (m *Mesh) UsesElementArray() bool {
	return m.usesElementArray
}
func (m *Mesh) VertexVBO() gl.Buffer {
	return m.vertexVBO
}
func (m *Mesh) VertexIndices() gl.Buffer {
	return m.vertexIndices
}
func (m *Mesh) NormalVBO() gl.Buffer {
	return m.normalVBO
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
func (m *Mesh) TextureCoords() gl.Buffer {
	return m.textureCoords
}

const (
	emptyTextureCoordsSize = 1024 * 1024
	emptyNormalsSize       = 1024 * 1024
	emptyIndicesSize       = 1024 * 1024
)

var (
	// TODO: Using these "empty" default buffers is ugly. I'm not sure how to do better though. If this is the way to go
	// make sure to document the size limitations as any meshes with undefined features that are larger than the size
	// limits will likely crash, or at least give a lot of error messages.

	// EmptyTexture is a texture buffer filled with just four bytes: [255, 255, 255, 255]
	// It is referenced by EmptyTextureCoords and is meant to be used as in meshes that don't contain another texture.
	EmptyTexture gl.Texture

	// EmptyTextureCoords contains emptyTextureCoordsSize floats worth of texture coordinates, all pointing to the single
	// texel in EmptyTexture. This allows it to be used as a stand in / default texture mapping for objects that don't
	// have a texture of their own.
	EmptyTextureCoords gl.Buffer

	// EmptyNormals contains emptyNormalsSize floats of zeros. This is used as a default for meshes that don't have
	// normals defined.
	EmptyNormals gl.Buffer

	EmptyIndices gl.Buffer
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
	coords := make([]float32, emptyTextureCoordsSize) // large array of 1 values
	for i := range coords {
		coords[i] = 1.0
	}
	textureCoordinates := bytecoder.Float32(binary.LittleEndian, coords...)

	EmptyTextureCoords = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, EmptyTextureCoords)
	gl.BufferData(gl.ARRAY_BUFFER, textureCoordinates, gl.STATIC_DRAW)
}

func initializeEmptyNormals() {
	normals := make([]float32, emptyNormalsSize) // large array of 0 values
	EmptyNormals = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, EmptyNormals)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, normals...), gl.STATIC_DRAW)
}

func initializeIndices() {
	indices := make([]float32, emptyIndicesSize) // large array of 0 values
	EmptyIndices = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, EmptyNormals)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, indices...), gl.STATIC_DRAW)
}
