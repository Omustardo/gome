package drawable

import (
	"image/color"

	"github.com/goxjs/gl"
)

type Drawable struct {
	Color   *color.RGBA
	Texture *gl.Texture
}
