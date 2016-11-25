package fps

import (
	"container/list"
	"time"
)

var _ tracker = (*rollingTracker)(nil)

// rollingFPSCounter keeps track of the calls to Update in the last second.
// It keeps a rolling window so it's completely accurate.
type rollingTracker struct {
	// framesLastSecond is a doubly linked list of the times of each call to rollingTracker.Update()
	// It has values greater than one second old removed regularly via rollingTracker.trim(),
	// so framesLastSecond.Len() should always be the exact number of frames in the last second.
	framesLastSecond *list.List

	// timeGetter allows mocking of time.Now() for testing.
	timeGetter func() time.Time
}

func newRollingCounter() *rollingTracker {
	return &rollingTracker{
		framesLastSecond: list.New(),
		timeGetter:       time.Now,
	}
}

// Update must be called once per frame.
func (f *rollingTracker) Update() {
	f.framesLastSecond.PushBack(f.timeGetter())
	f.trim()
}

// trim removes any recorded timestamps that are over a second old. Since new timestamps are only added once per frame,
// and rollingTracker should be called once per frame, the expected cost of this function is only checking and
// removing a single element.
func (f *rollingTracker) trim() {
	now := f.timeGetter()
	for e := f.framesLastSecond.Front(); e != nil; {
		timeToCheck := e.Value.(time.Time)
		if now.Sub(timeToCheck) > time.Second {
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
func (f *rollingTracker) DeltaTime() time.Duration {
	// This should only happen when the game first loads, and only for two frames, so shouldn't be noticeable by the user.
	if f.framesLastSecond.Len() < 2 {
		return time.Nanosecond * 0
	}
	lastCall := f.framesLastSecond.Back().Value.(time.Time)
	prevLastCall := f.framesLastSecond.Back().Prev().Value.(time.Time)
	return lastCall.Sub(prevLastCall)
}

// DeltaTime returns the time elapsed between the last two calls to Update, in seconds.
func (f *rollingTracker) DeltaTimeSeconds() float32 {
	return float32(f.DeltaTime().Seconds())
}
