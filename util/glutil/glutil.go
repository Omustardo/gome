// glutil contains common functions to deal with OpenGl and WebGL, such as loading data into Buffers and Textures.
package glutil

import (
	"fmt"

	"github.com/goxjs/gl"
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
