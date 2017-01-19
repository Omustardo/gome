package asset

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"

	// for decoding of different file types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/goxjs/gl"
)

// LoadTexture from local assets. Handles jpg, png, and static gifs.
func LoadTexture(path string) (gl.Texture, error) {
	fileData, err := loadFile(path)
	if err != nil {
		return gl.Texture{}, err
	}

	img, _, err := image.Decode(bytes.NewBuffer(fileData))
	if err != nil {
		return gl.Texture{}, err
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Get raw RGBA pixels by drawing the image into an NRGBA image. This is necessary to deal with different
	// image formats that aren't decoded into a pixel array. For example, jpeg compressed images are read in in a way
	// that mimics their encoding, and due to the way they are compressed, you can't get pixel values easily.
	// By drawing the decoded image out as NRGBA, we are guaranteed to get something we can deal with.
	// Note that this is wasteful for images that are already read in as RGBA or NRGBA, but it's a one time cost
	// and shouldn't be an issue.
	newimg := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(newimg, bounds, img, bounds.Min, draw.Src)
	data := newimg.Pix

	// Need to flip the image vertically since OpenGL considers 0,0 to be the top left corner.
	// Note width*4 since the data array consists of R,G,B,A values.
	if err := flipYCoords(data, width*4); err != nil {
		return gl.Texture{}, err
	}
	return LoadTextureData(width, height, data), nil
}

// LoadTextureData takes raw RGBA image data and puts it into a texture unit on the GPU.
// It's up to the caller to delete the texture buffer using gl.DeleteTexture(texture) when it's no longer needed.
func LoadTextureData(width, height int, data []uint8) gl.Texture {
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
	return texture
}

// Takes a flattened 2D array and the width of the rows.
// Modifies the values such that if the original array were an image, it would now appear upside down.
func flipYCoords(data []uint8, width int) error {
	if len(data)%width != 0 {
		return fmt.Errorf("expected flattened 2d array, got uneven row length: len %% width == %v", len(data)%width)
	}
	height := len(data) / width
	for row := 0; row < height/2; row++ {
		for col := 0; col < width; col++ {
			temp := data[col+row*width]
			data[col+row*width] = data[col+(height-1-row)*width]
			data[col+(height-1-row)*width] = temp
		}
	}
	return nil
}
