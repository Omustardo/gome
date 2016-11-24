// Package fps can be used to keep track of framerate.
// Sample usage:
//
package fps

import "github.com/omustardo/gome/util"

// FPS Counter keeps track of the number of frames in the last second.
type FPSCounter struct {
	framesLastSecond int
	prevTime         int64
	frameCounter     int
}

func NewFPSCounter() *FPSCounter {
	return &FPSCounter{
		prevTime: util.GetTimeMillis(),
	}
}

// Update must be called once per frame.
func (f *FPSCounter) Update() {
	// TODO: Improve tracker so it actually keeps track of the last second's worth of events.
	currTime := util.GetTimeMillis()
	f.frameCounter++
	if currTime-f.prevTime >= 1000 {
		f.framesLastSecond = f.frameCounter
		f.frameCounter = 0
		f.prevTime = currTime
	}
}

// GetFPS returns the number of frames shown in the last second.
// Note that the counter is updated every second, so at most this could be the time period 0.999 to 1.999 seconds ago,
// rather than the past second.
func (f *FPSCounter) GetFPS() int {
	return f.framesLastSecond
}

// GetFramerate returns an estimate of the average length of each frame.
func (f *FPSCounter) GetFramerate() float32 {
	return 1000 / float32(f.framesLastSecond)
}
