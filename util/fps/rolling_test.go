package fps

import (
	"container/list"
	"testing"
	"time"
)

type MockTime struct {
	timesMillis []int64
}

func (m *MockTime) Now() time.Time {
	if len(m.timesMillis) < 0 {
		panic("no times left to provide")
	}
	timeToReturn := m.timesMillis[0]
	if len(m.timesMillis) == 1 {
		m.timesMillis = nil
	} else {
		m.timesMillis = m.timesMillis[1:]
	}
	return time.Unix(timeToReturn/1000, 1e6*(timeToReturn%1000))
}

func TestRollingTracker(t *testing.T) {
	tests := []struct {
		mockTime    *MockTime
		wantDelta   time.Duration
		wantFPS     int
		updateCount int
	}{
		{
			mockTime: &MockTime{
				timesMillis: []int64{
					1, 1, // Need to use each value twice because the mock gets called twice per update. Once in Update(), once in trim().
					16, 16,
					38, 38,
					55, 55,
					55, // FPS() calls trim(), which gets the time, so provide an extra value.
				},
			},
			updateCount: 4, // calls to update, which consume 2*updateCount of the mockTime values
			wantDelta:   time.Millisecond * 17,
			wantFPS:     4,
		},
		{
			mockTime: &MockTime{
				timesMillis: []int64{
					1, 1,
					500, 500,
					999, 999,
					1000, 1000,
					1500, 1500,
					2000, 2000,
					2000,
				},
			},
			updateCount: 6,
			wantDelta:   time.Millisecond * 500,
			wantFPS:     3, // The FPS() call at time 2000 includes the frame from time 1000, but none before that.
		},
	}

	for _, tt := range tests {
		fpsTracker := rollingTracker{
			framesLastSecond: list.New(),
			timeGetter:       tt.mockTime.Now,
		}
		for i := 0; i < tt.updateCount; i++ {
			fpsTracker.Update()
		}

		gotDelta := fpsTracker.DeltaTime()
		if gotDelta != tt.wantDelta {
			t.Errorf("Want: %v delta time, got: %v", tt.wantDelta, gotDelta)
		}
		gotFPS := fpsTracker.FPS()
		if gotFPS != tt.wantFPS {
			t.Errorf("Want: %v fps, got: %v", tt.wantFPS, gotFPS)
		}
	}
}
