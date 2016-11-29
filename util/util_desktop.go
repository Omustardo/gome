// +build !js

package util

import (
	"image"
	"image/png"
	"os"

	"github.com/goxjs/gl"
)

// SaveScreenshot reads pixel data from OpenGL buffers, so it must be run in the same main thread as the rest
// of OpenGL.
// TODO: write to file in a goroutine and return a (chan err), or just ignore slow errors. Handling errors that can be caught immediately is fine. Blocking while writing to file adds way too much delay.
func SaveScreenshot(width, height int, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	gl.ReadPixels(img.Pix, 0, 0, width, height, gl.RGBA, gl.UNSIGNED_BYTE)

	// Need to flip the image vertically since the pixels are provided with (0,0) in the top left corner.
	FlipImageVertically(img)

	out, err := os.Create(path) // TODO: WebGL isn't happy with this (no syscalls allowed). Need to make a util_js.go with conditional compilation.
	if err != nil {
		return err
	}
	return png.Encode(out, img)
}
