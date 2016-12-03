// zoom is a handler for dealing with zoom.
package zoom

import (
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
)

// arbitrary number of scroll wheel ticks to change zoom from min to max. Balanced between large size for quick scrolling, and small size for smooth zooming. TODO: Use larger size, but interpolate so it's smooth no matter what.
const step = 60

type Zoom interface {
	// GetCurrentPercent returns the current percent zoom. Always a positive value.
	GetCurrentPercent() float32

	// Range returns the limits of zooming in and out.
	// min is the limit of zooming out, max is the limit of zooming in.
	// For example a min of 0.25 means objects appear to be 25% of their original size when zoomed out fully.
	// A max of 3 means objects will appear to be 300% of their original size when fully zoomed in.
	Range() (min, max float32)

	// Update must be called every frame.
	Update()
}

// ScrollZoom implements Zoom. Intended for use with a camera's projection matrix. Intended to get data from a mouse scroll wheel.
type ScrollZoom struct {
	// Range of percent zoom allowed. If the range doesn't include 1.0, then the default starting zoom will be (Min+Max)/2.
	// Min zoom means zoomed out as far as possible - everything will look small. Max zoom is zoomed in as close as possible.
	Min, Max float32
	// curr is the current percent. It is updated in the Update() function.
	curr float32
	// GetScrollAmount is expected to return the amount scrolled in the last update period.
	GetScrollAmount func() float32
	// Percent to zoom in per scroll wheel "tick".
	percentPerScrollAmount float32
}

// NewScrollZoom creates a ScrollZoom struct.
// Example usage:
// 	zoomer := zoom.NewScrollZoom(0.25, 3, // allows zooming out to 25% of the original size, and in to 300% of the original size.
//	  func() float32 { return mouse.Handler.Scroll().Y() },
//	)
func NewScrollZoom(min, max float32, GetScrollAmount func() float32) *ScrollZoom {
	if min <= 0 {
		panic(fmt.Sprintf("invalid min zoom: %v < 0", min))
	}
	if GetScrollAmount == nil {
		panic("GetScrollAmount is undefined")
	}
	// Try to default to no zoom. If range doesn't allow it, use the average.
	curr := float32(1.0)
	if min > 1 || max < 1 {
		curr = (min + max) / 2
	}
	return &ScrollZoom{
		Min:                    min,
		Max:                    max,
		curr:                   curr,
		GetScrollAmount:        GetScrollAmount,
		percentPerScrollAmount: (max - min) / step,
	}
}

func (z *ScrollZoom) Update() {
	if ticks := z.GetScrollAmount(); ticks != 0 {
		percentChange := z.percentPerScrollAmount * ticks
		z.curr = mgl32.Clamp(z.curr+percentChange, z.Min, z.Max)
	}
}

func (z *ScrollZoom) Range() (min, max float32) {
	return z.Min, z.Max
}

func (z *ScrollZoom) GetCurrentPercent() float32 {
	return z.curr
}
