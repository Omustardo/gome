// Package fps can be used to keep track of framerate.
// Sample usage:
//   fps.Initialize()
//   for { // game loop
//     fps.Handler.Update()
// 	   // do game stuff
//     fmt.Println(fps.GetFPS())
//   }
package fps

import (
	"time"
)

var Handler *fpsCounter

// fpsCounter keeps track of the number of frames in the last second.
type fpsCounter struct {
	framesLastSecond int
	prevTime         time.Time
	frameCounter     int
}

func Initialize() {
	if Handler != nil {
		panic("fps.Handler already initialized")
	}
	Handler = &fpsCounter{
		prevTime: time.Now(),
	}
}

// Update must be called once per frame.
func (f *fpsCounter) Update() {
	// TODO: Improve tracker so it actually keeps track of the last second's worth of events with a rolling window, rather than discrete one second windows.
	currTime := time.Now()
	f.frameCounter++
	if currTime.Sub(f.prevTime) >= time.Second {
		f.framesLastSecond = f.frameCounter
		f.frameCounter = 0
		f.prevTime = currTime
	}
}

// GetFPS returns the number of frames shown in the last second.
// Note that the counter is updated every second, so at most this could be the time period 0.999 to 1.999 seconds ago,
// rather than the past second.
func (f *fpsCounter) GetFPS() int {
	return f.framesLastSecond
}

// GetFramerate returns an estimate of the average length of each frame.
func (f *fpsCounter) GetFramerate() float32 {
	if f.framesLastSecond == 0 {
		return -1
	}
	return 1000 / float32(f.framesLastSecond)
}
