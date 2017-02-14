package asset

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"log"

	// for decoding of different file types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util"
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
	return LoadTextureDataImageNRGBA(newimg), nil
}

// LoadTextureData takes raw RGBA image data and puts it into a texture unit on the GPU.
// It's up to the caller to delete the texture buffer using gl.DeleteTexture(texture) when it's no longer needed.
// Note that the input data must not be ragged - each row must be the same length.
func LoadTextureData2D(data [][]uint8) gl.Texture {
	if len(data) == 0 {
		return gl.Texture{}
	}
	width := len(data[0])
	flat := make([]uint8, len(data)*len(data[0]))
	for _, row := range data {
		if len(row) != width {
			log.Printf("Got ragged 2D array. Found rows of len %d and %d", width, len(row))
			return gl.Texture{}
		}
		flat = append(flat, row...)
	}
	return LoadTextureData(width, flat)
}

func LoadTextureDataImageRGBA(img *image.RGBA) gl.Texture {
	if img == nil {
		return gl.Texture{}
	}
	return LoadTextureData(img.Bounds().Dx(), img.Pix)
}
func LoadTextureDataImageNRGBA(img *image.NRGBA) gl.Texture {
	if img == nil {
		return gl.Texture{}
	}
	return LoadTextureData(img.Bounds().Dx(), img.Pix)
}

// LoadTextureData takes raw RGBA image data and puts it into a texture unit on the GPU.
// width is the length of each row of data in RGBA pixels - note that the input data is in bytes, so the actual length
// of each row is width * 4, and the total length of the input data should be 4 * width * height. TODO: Consider passing in actual width of the data. This is a bit confusing either way... alternatively pass in array of colors, but that loses efficiency.
//
// It's up to the caller to delete the texture buffer using gl.DeleteTexture(texture) when it's no longer needed.
// Note that the input data must not be ragged - each row must be the same length.
func LoadTextureData(width int, data []uint8) gl.Texture {
	if width <= 0 {
		return gl.Texture{} // TODO: return errors
	}
	if !util.IsPowerOfTwo(width) || !util.IsPowerOfTwo(len(data)/width) {
		return gl.Texture{} // TODO: return errors: webgl requires poweroftwo texture dimensions to make mipmaps.
	}

	// gl.Enable(gl.TEXTURE_2D) // some sources says this is needed, but it doesn't seem to be. In fact, it gives an "invalid capability" message in webgl.
	height := len(data) / (width * 4)

	// TODO: Am I doing textures incorrectly? The freetype demos don't flip and the image still comes out fine. I think that is just that image.Encode uses 0,0 as bottom left, but double check to be sure.
	// Need to flip the image vertically since OpenGL considers 0,0 to be the top left corner.
	if err := flipYCoords(width*4, data); err != nil {
		return gl.Texture{}
	}

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
func flipYCoords(width int, data []uint8) error {
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
