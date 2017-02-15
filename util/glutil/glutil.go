// glutil contains common functions to deal with OpenGl and WebGL, such as loading data into Buffers and Textures.
package glutil

import (
	"encoding/binary"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/bytecoder"
	"github.com/omustardo/gome/util"
)

// LoadTextureData takes raw RGBA image data and puts it into a texture unit on the GPU.
// Width and height are refer to the actual image, not the number of RGBA bytes. For example, an 8 by 8 pixel image
// should be provided as (8, 8, data) even though the size of the data should be 8 * 8 * 4, since there are 4 bytes per
// pixel.
//
// Errors are returned if:
// - data is the incorrect size for the provided width and height
// - either width or height is not a power of 2. This is a requirement for making mipmaps in webgl.
//
// TODO: Example usage. Make sure to give an example that uses image.Pix as data.
//
// It's up to the caller to delete the texture buffer using gl.DeleteTexture(texture) when it's no longer needed.
// Note that the input data must not be ragged - each row must be the same length.
func LoadTextureData(width, height int, data []uint8) (gl.Texture, error) {
	// Note that we can expect height to equal: len(data) / (width * 4)
	// but asking users to pass it in is easier and avoids edge cases, like if the input is invalid.

	if len(data)%4 != 0 {
		return gl.Texture{}, fmt.Errorf("data length must be a multiple 4 to represent RGBA pixels, got: %d", len(data)/4)
	}
	if width*height*4 != len(data) {
		return gl.Texture{}, fmt.Errorf("image dimensions don't match input data length: %d * %d != %d", width, height, len(data)/4)
	}
	if width <= 0 || height <= 0 {
		return gl.Texture{}, fmt.Errorf("image dimensions must be >0. got [%d,%d]", width, height)
	}
	if !util.IsPowerOfTwo(width) || !util.IsPowerOfTwo(len(data)/width) {
		return gl.Texture{}, fmt.Errorf("image dimensions must be powers of two. got [%d,%d]", width, height)
	}

	// gl.Enable(gl.TEXTURE_2D) // some sources says this is needed, but it doesn't seem to be. In fact, it gives an "invalid capability" message in webgl.

	texture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// NOTE: gl.FLOAT isn't enabled for texture data types unless gl.getExtension('OES_texture_float'); is set, so just use gl.UNSIGNED_BYTE
	//   See http://stackoverflow.com/questions/23124597/storing-floats-in-a-texture-in-opengl-es  http://stackoverflow.com/questions/22666556/webgl-texture-creation-trouble
	gl.TexImage2D(gl.TEXTURE_2D, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, gl.Texture{}) // bind to "null" to prevent using the wrong texture by mistake.
	return texture, nil
}

// LoadBuffer takes a float32 slice and stores the underlying data in a buffer on the GPU.
func LoadBuffer(data []byte) gl.Buffer {
	buf := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)                  // Bind the target buffer so we can store values in it. https://www.opengl.org/sdk/docs/man4/html/glBindBuffer.xhtml
	gl.BufferData(gl.ARRAY_BUFFER, data, gl.STATIC_DRAW) // store values in buffer
	return buf
}

// LoadBufferUint16 takes a uint16 slice and stores the underlying data in a buffer on the GPU.
func LoadBufferUint16(data []uint16) gl.Buffer {
	return LoadBuffer(bytecoder.Uint16(binary.LittleEndian, data...))
}

// LoadBufferFloat32 takes a float32 slice and stores the underlying data in a buffer on the GPU.
func LoadBufferFloat32(data []float32) gl.Buffer {
	return LoadBuffer(bytecoder.Float32(binary.LittleEndian, data...))
}

func LoadBufferVec2(data []mgl32.Vec2) gl.Buffer {
	return LoadBuffer(bytecoder.Vec2(binary.LittleEndian, data...))
}

func LoadBufferVec3(data []mgl32.Vec3) gl.Buffer {
	return LoadBuffer(bytecoder.Vec3(binary.LittleEndian, data...))
}

// LoadIndexBuffer takes a uint16 slice and stores the underlying data in an ELEMENT_ARRAY buffer on the GPU.
func LoadIndexBuffer(data []uint16) gl.Buffer {
	buf := gl.CreateBuffer()
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, data...), gl.STATIC_DRAW)
	return buf
}
