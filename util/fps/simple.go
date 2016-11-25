package fps

import "time"

var _ tracker = (*simpleTracker)(nil)

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

	// timeGetter allows mocking of time.Now() for testing.
	timeGetter func() time.Time
}

func newSimpleTracker() tracker {
	return &simpleTracker{
		prevTime:   time.Now(),
		timeGetter: time.Now,
	}
}

// Update must be called once per frame.
func (f *simpleTracker) Update() {
	currTime := f.timeGetter()
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
func (f *simpleTracker) DeltaTime() time.Duration {
	return f.lastUpdateTime.Sub(f.prevLastUpdateTime)
}

// DeltaTime returns the time elapsed between the last two calls to Update, in seconds.
func (f *simpleTracker) DeltaTimeSeconds() float32 {
	return float32(f.DeltaTime().Seconds())
}
