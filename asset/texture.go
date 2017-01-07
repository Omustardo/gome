package asset

import (
	"bytes"
	"fmt"
	"image"

	// for decoding of different file types
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"image/color"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util"
)

// LoadTexture from local assets.
func LoadTexture(path string) (gl.Texture, error) {
	// based on https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial/Using_textures_in_WebGL and https://golang.org/pkg/image/

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

	// Need to flip the image vertically since OpenGL considers 0,0 to be the top left corner.
	util.FlipImageVertically(img)

	// Image checking from https://github.com/go-gl-legacy/glh/blob/master/texture.go
	switch trueim := img.(type) {
	case *image.RGBA:
		return LoadTextureData(width, height, trueim.Pix), nil
	case *image.NRGBA: // NRGBA is non-premultiplied RGBA. RGBA evidently is supposed to multiply the alpha value by the other colors, so NRGBA of (1, 0.5, 0, 0.5) is RGBA of (0.5, 0.25, 0. 0.5)
		return LoadTextureData(width, height, trueim.Pix), nil
	case *image.YCbCr:
		return LoadTextureData(width, height, ycbCrToRGBA(trueim).Pix), nil
	default:
		// copy := image.NewRGBA(trueim.Bounds())
		// draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		return gl.Texture{}, fmt.Errorf("unsupported texture format %T", img)
	}
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

func ycbCrToRGBA(img *image.YCbCr) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			c := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			rgba.SetRGBA(x, y, c)
		}
	}
	return rgba
}
