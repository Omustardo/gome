package drawable

import (
	"image/color"

	"github.com/goxjs/gl"
)

type Attributes struct {
	Color   *color.NRGBA
	Texture *gl.Texture
}

type Drawable interface {
	// VBO() *gl.Buffer
	// Attributes() Attributes
	// Entity() entity.Entity
}