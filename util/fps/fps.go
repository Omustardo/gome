// Package fps can be used to keep track of framerate and delta time.
// Sample usage:
//   fps.Initialize()
//   for { // game loop
//     fps.Handler.Update()
//     fmt.Println(fps.GetFPS())
// 	   // do game stuff
//   }
package fps

import (
	"time"
)

// Handler is the singleton fps handler. It should be initialized with fps.Initialize(), and then
// all framerate and delta time checking can access it directly.
var Handler tracker

func Initialize() {
	if Handler != nil {
		panic("Handler already initialized")
	}
	Handler = newRollingCounter()
}

type tracker interface {
	// Update must be called once per frame.
	Update()
	// FPS returns the number of frames shown in the last second.
	FPS() int
	// Framerate returns an estimate of the average millisecond duration of each frame. i.e. 16.66
	Framerate() float32
	// DeltaTime returns the time elapsed between the last two calls to Update.
	DeltaTime() time.Duration
	// DeltaTimeSeconds returns the time elapsed between the last two calls to Update, in seconds.
	DeltaTimeSeconds() float32
}
