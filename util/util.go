package util

import (
	"image"
	"time"

	"fmt"
	"math/rand"
)

// Why is this not in the standard time library? Am I missing something?
func GetTimeMillis() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func FlipImageVertically(img image.Image) error {
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	// Image checking from https://github.com/go-gl-legacy/glh/blob/master/texture.go
	switch trueim := img.(type) {
	case *image.RGBA:
		for row := 0; row < height/2; row++ {
			for col := 0; col < width; col++ {
				temp := img.At(col, row)
				trueim.Set(col, row, img.At(col, height-1-row))
				trueim.Set(col, height-1-row, temp)
			}
		}
	case *image.NRGBA: // What is NRGBA? It seems to act exactly like RGBA.
		for row := 0; row < height/2; row++ {
			for col := 0; col < width; col++ {
				temp := img.At(col, row)
				trueim.Set(col, row, img.At(col, height-1-row))
				trueim.Set(col, height-1-row, temp)
			}
		}
	default:
		return fmt.Errorf("unknown image type: %T", img)
	}
	return nil
}

func RandUint8() uint8 {
	return uint8(rand.Uint32() % 256)
}
