// Package asset provides tools to manage loading, use, and unloading of assets, such as images and audio.
package asset

import (
	"bytes"
	"fmt"
	"image"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util"
)

// LoadTexture from local assets.
func LoadTexture(path string) (*gl.Texture, error) {
	// based on https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Tutorial/Using_textures_in_WebGL and https://golang.org/pkg/image/

	fileData, err := loadFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewBuffer(fileData))
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Need to flip the image vertically since OpenGL considers 0,0 to be the top left corner.
	util.FlipImageVertically(img)

	// Image checking from https://github.com/go-gl-legacy/glh/blob/master/texture.go
	switch trueim := img.(type) {
	case *image.RGBA:
		return LoadTextureData(width, height, trueim.Pix), nil
	case *image.NRGBA: // What is NRGBA? It seems to act exactly like RGBA.
		return LoadTextureData(width, height, trueim.Pix), nil
	default:
		// copy := image.NewRGBA(trueim.Bounds())
		// draw.Draw(copy, trueim.Bounds(), trueim, image.Pt(0, 0), draw.Src)
		return nil, fmt.Errorf("unsupported texture format %T", img)
	}
}

func LoadTextureData(width, height int, data []uint8) *gl.Texture {
	// gl.Enable(gl.TEXTURE_2D) // some sources says this is needed, but it doesn't seem to be. In fact, it gives an "invalid capability" message in webgl.

	texture := gl.CreateTexture()
	gl.BindTexture(gl.TEXTURE_2D, texture)
	// NOTE: gl.FLOAT isn't enabled for texture data types unless gl.getExtension('OES_texture_float'); is set, so just use gl.UNSIGNED_BYTE
	//   See http://stackoverflow.com/questions/23124597/storing-floats-in-a-texture-in-opengl-es  http://stackoverflow.com/questions/22666556/webgl-texture-creation-trouble
	gl.TexImage2D(gl.TEXTURE_2D, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE, data) // TODO: Does layering RGBA images work? Or do we need to sort by Z value and draw in that order.
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.GenerateMipmap(gl.TEXTURE_2D)
	// gl.BindTexture(gl.TEXTURE_2D, gl.Texture{Value: 0}) // in js demo, they bind to null to prevent using the wrong texture by mistake. No way to do that with structs?
	return &texture
}
