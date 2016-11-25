// Package fps can be used to keep track of framerate.
// Sample usage:
//   fps.Initialize()
//   for { // game loop
//     fps.Handler.Update()
//     fmt.Println(fps.GetFPS())
// 	   // do game stuff
//   }
package fps

import (
	"container/list"
	"log"
	"time"
)

var Handler tracker

func Initialize() {
	if Handler != nil {
		panic("fps.Handler already initialized")
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
	// DeltaTime returns the millisecond time elapsed between the last two calls to Update.
	DeltaTime() float32
}

// rollingFPSCounter keeps track of the number of frames in the last second.
// It keeps a rolling window so it's completely accurate.
type rollingTracker struct {
	// framesLastSecond is a doubly linked list of the times of each call to rollingTracker.Update()
	// It has values greater than one second old removed regularly via rollingTracker.trim(),
	// so framesLastSecond.Len() should always be the exact number of frames in the last second.
	framesLastSecond *list.List
}

func newRollingCounter() *rollingTracker {
	return &rollingTracker{
		framesLastSecond: list.New(),
	}
}

// Update must be called once per frame.
func (f *rollingTracker) Update() {
	f.framesLastSecond.PushBack(time.Now())
	f.trim()
}

// trim removes any recorded timestamps that are over a second old. Since new timestamps are only added once per frame,
// and rollingTracker should be called once per frame, the expected cost of this function is only checking and
// removing a single element.
func (f *rollingTracker) trim() {
	for e := f.framesLastSecond.Front(); e != nil; {
		if timeToCheck, ok := e.Value.(time.Time); !ok {
			log.Printf("unexpected type %T. expected time.Time", e.Value)
		} else if time.Now().Sub(timeToCheck) > time.Second {
			temp := e
			e = e.Next()
			f.framesLastSecond.Remove(temp)
		} else {
			// Since timestamps are provided in increasing order, once one's large enough, ignore the rest.
			return
		}
	}
}

// FPS returns the number of frames shown in the last second.
func (f *rollingTracker) FPS() int {
	f.trim()
	return f.framesLastSecond.Len()
}

// Framerate returns an estimate of the average duration of each frame.
func (f *rollingTracker) Framerate() float32 {
	if f.framesLastSecond.Len() == 0 {
		return -1
	}
	return 1000 / float32(f.framesLastSecond.Len())
}

// DeltaTime returns the time elapsed between the last two calls to Update.
func (f *rollingTracker) DeltaTime() float32 {
	// This should only happen when the game first loads, and only for two frames, so shouldn't be noticeable by the user.
	if f.framesLastSecond.Len() < 2 {
		return 0
	}
	lastCall := f.framesLastSecond.Back().Value.(time.Time)
	prevLastCall := f.framesLastSecond.Back().Prev().Value.(time.Time)
	return float32(lastCall.Sub(prevLastCall).Nanoseconds()) / float32(time.Millisecond/time.Nanosecond)
}

// simpleTracker keeps track of the number of frames in the last second.
// It is extremely efficient in both time and space, but it only keeps track of the average frame rate
// in the last full second. It updates with discrete one second windows, so at worst it could be the time period 0.999
// to 1.999 seconds ago, rather than the exact last second. It is recommended to only use this if you target extremely
// low spec hardware.
type simpleTracker struct {
	framesLastSecond   int
	prevTime           time.Time
	lastUpdateTime     time.Time
	prevLastUpdateTime time.Time
	frameCounter       int
}

func newSimpleTracker() tracker {
	return &simpleTracker{
		prevTime: time.Now(),
	}
}

// Update must be called once per frame.
func (f *simpleTracker) Update() {
	currTime := time.Now()
	f.prevLastUpdateTime = f.lastUpdateTime
	f.lastUpdateTime = currTime

	f.frameCounter++
	if currTime.Sub(f.prevTime) > time.Second {
		f.framesLastSecond = f.frameCounter
		f.frameCounter = 0
		f.prevTime = currTime
	}
}

// FPS returns the number of frames shown in the last second.
func (f *simpleTracker) FPS() int {
	return f.framesLastSecond
}

// Framerate returns an estimate of the average duration of each frame.
func (f *simpleTracker) Framerate() float32 {
	if f.framesLastSecond == 0 {
		return -1
	}
	return 1000 / float32(f.framesLastSecond)
}

// DeltaTime returns the time elapsed between the last two calls to Update.
func (f *simpleTracker) DeltaTime() float32 {
	return float32(f.lastUpdateTime.Sub(f.prevLastUpdateTime).Nanoseconds()) / float32(time.Millisecond/time.Nanosecond)
}
