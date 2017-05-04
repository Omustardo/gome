package zoom_test

import (
	"github.com/omustardo/gome/camera/zoom"
	"log"
	"testing"
)

func TestScrollZoom(t *testing.T) {
	tests := []struct {
		mouseY, zoomPercent float32
	} {
		{
			0, 1.0,
		},
		{
			30, 2.375,
		},
		{
			-30, 0.25,
		},
		{
			60, 3.0,
		},
	}

	for _, tt := range tests {
		zoomer := zoom.NewScrollZoom(
			// Allow zooming out to 25% of the original size, and in to 300% of the original size.
			0.25, 3,
			// For this example, simulate getting mouse input by changing this variable.
			func() float32 { return tt.mouseY },
		)
		zoomer.Update()

		if current := zoomer.GetCurrentPercent(); current != tt.zoomPercent  {
			t.Errorf("for mouse scroll %v, got current zoom percent = %v, expected %v", tt.mouseY, current, tt.zoomPercent)
		}
	}
}

func Example() {
	var mouseScrollY float32

	zoomer := zoom.NewScrollZoom(
		// Allow zooming out to 25% of the original size, and in to 300% of the original size.
		0.25, 3,
		// For this example, simulate getting mouse input by changing this variable.
		func() float32 { return mouseScrollY },
	)

	log.Println(zoomer.GetCurrentPercent()) // 1.0 : no change since the mouse hasn't scrolled yet.

	mouseScrollY = 30
	log.Println(zoomer.GetCurrentPercent()) // value > 1.0 : the mouse was scrolled.
}