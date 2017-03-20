package asset

import (
	"bytes"
	"errors"
	"image"
	"image/draw"

	// for decoding of different file types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/glutil"
)

// LoadTexture from local assets. Handles jpg, png, and static gifs.
func LoadTexture(path string) (gl.Texture, error) {
	fileData, err := LoadFile(path)
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
	return LoadTextureDataImage(img)
}

// LoadTextureData takes raw RGBA image data and puts it into a texture unit on the GPU.
// It's up to the caller to delete the texture buffer using gl.DeleteTexture(texture) when it's no longer needed.
// Note that the input data must not be ragged - each row must be the same length.
func LoadTextureData2D(data [][]uint8) (gl.Texture, error) {
	if len(data) == 0 {
		return gl.Texture{}, errors.New("no data provided")
	}
	width, height := len(data[0]), len(data)
	flat := flatten(data)
	flipVertically(width, flat)
	return glutil.LoadTextureData(width, height, flat)
}

// LoadTextureDataImage loads the provided image onto the GPU and returns a reference to the gl Texture.
func LoadTextureDataImage(img image.Image) (gl.Texture, error) {
	if img == nil {
		return gl.Texture{}, errors.New("nil image can't be loaded")
	}
	// TODO: Am I doing textures incorrectly? The freetype demos don't flip and the image still comes out fine. I think that is just that image.Encode uses 0,0 as bottom left, but double check to be sure.
	// Need to flip the image vertically since OpenGL considers 0,0 to be the top left corner.
	imgNRGBA := flipImageVertically(img)
	return glutil.LoadTextureData(imgNRGBA.Bounds().Dx(), imgNRGBA.Bounds().Dy(), imgNRGBA.Pix)
}

func flipImageVertically(img image.Image) *image.NRGBA {
	if img == nil {
		return nil
	}
	newimg := image.NewNRGBA(img.Bounds())
	for row := 0; row < img.Bounds().Dy(); row++ {
		for col := 0; col < img.Bounds().Dx(); col++ {
			newimg.Set(col, img.Bounds().Dy()-1-row, img.At(col, row))
		}
	}
	return newimg
}

// flipVertically takes a flattened 2D array and the width of the rows.
// Modifies the values such that if the original array were an image, it would now appear upside down.
// This is a best attempt. If the provided array has ragged rows when converted to 2D, the invalid locations are ignored.
func flipVertically(width int, data []uint8) {
	height := len(data) / width
	for row := 0; row < height/2; row++ {
		for col := 0; col < width; col++ {
			if len(data) < (col+row*width) && len(data) < col+(height-1-row)*width {
				temp := data[col+row*width]
				data[col+row*width] = data[col+(height-1-row)*width]
				data[col+(height-1-row)*width] = temp
			}
		}
	}
}

func flatten(data [][]uint8) []uint8 {
	var flat []uint8
	for i := range data {
		flat = append(flat, data[i]...)
	}
	return flat
}
