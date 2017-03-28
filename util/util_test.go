package util

import "testing"

func TestIsPowerOfTwoValid(t *testing.T) {
	n := 1

	// Loop until overflow. The number of bytes in an int depends on architecture of the system,
	// so the loop can't be hardcoded.
	for n > 0 {
		if !IsPowerOfTwo(n) {
			t.Errorf("got %v is not a power of two", n)
		}
		n <<= 1
	}
}

func TestIsPowerOfTwoInvalid(t *testing.T) {
	tests := []int{-2, -1, 0, 3}

	for _, tt := range tests {
		if IsPowerOfTwo(tt) {
			t.Errorf("incorrectly got %v is a power of two", tt)
		}
	}
}

func TestRoundUpToPowerOfTwo(t *testing.T) {
	tests := []struct {
		n, want int
	}{
		{
			n:    0,
			want: 1,
		},
		{
			n:    1,
			want: 1,
		},
		{
			n:    2,
			want: 2,
		},
		{
			n:    3,
			want: 4,
		},
		{
			n:    -1,
			want: 1,
		},
		{
			n:    -1000,
			want: 1,
		},
	}

	for _, tt := range tests {
		if got := RoundUpToPowerOfTwo(tt.n); got != tt.want {
			t.Errorf("RoundUpToPowerOfTwo(%d)=%d; wanted %d", tt.n, got, tt.want)
		}
	}
}
