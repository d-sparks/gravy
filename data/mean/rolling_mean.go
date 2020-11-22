package mean

import (
	"math"
)

// Rolling is a struct that tracks n-day averages for an asset or portfolio.
type Rolling struct {
	buffer []float64
	sums   map[int]float64
	head   int
	tails  map[int]int
	filled int
}

// New creates a moving average struct that tracks the given windows.
func NewRolling(days ...int) *Rolling {
	rolling := Rolling{sums: map[int]float64{}, head: 0, tails: map[int]int{}, filled: 0}
	maxRequest := 0
	for _, n := range days {
		if n <= 0 {
			panic("Can't track 0 day average.")
		}
		if n > maxRequest {
			maxRequest = n
		}
		rolling.sums[n] = 0.0
		rolling.tails[n] = 0
	}
	rolling.buffer = make([]float64, maxRequest)
	return &rolling
}

// Observe inserts the value. The way this works is there is one large cyclic buffer whose size is equal to the longest
// moving average being tracked, however there are several tails where the shorter moving averages end.
func (r *Rolling) Observe(value float64) {
	// Update the sums and tails if necessary.
	for n := range r.sums {
		if r.filled >= n {
			r.sums[n] -= r.buffer[r.tails[n]]
			r.tails[n] = (r.tails[n] + 1) % len(r.buffer)
		}
		r.sums[n] += value
	}

	// Insert into cyclic buffer and advance the head index.
	r.buffer[r.head] = value
	r.head = (r.head + 1) % len(r.buffer)
	if r.filled < len(r.buffer) {
		r.filled++
	}
}

// Value returns the moving average for the prescribed timeframe. Begins tracking the timeframe if it is not already
// being tracked. Assumes n is > 0 and that there has been at least one observation.
func (r *Rolling) Value(n int) float64 {
	if n <= 0 {
		panic("Can't track 0 day average.")
	}

	// Return the average if it is already tracked.
	if sum, ok := r.sums[n]; ok {
		return sum / math.Min(float64(n), float64(r.filled))
	}

	// Otherwise, begin tracking. The tail will be at r.head - n % len(r.buffer) unless the buffer is not yet big
	// enough or not yet full enough. When the buffer is not full enough, e.g. n >= r.filled, we want the tail to be
	// at 0, or r.head - r.filled. When the buffer is not big enough, we want to start at the oldest element, which
	// is at r.head, i.e. we want r.head - len(r.buffer) = r.head. This can all be expressed as
	//
	//   r.tails[n] = (r.head + len(r.buffer) - min(n, r.filled, len(r.buffer)) % len(r.buffer)
	//
	// However, since Go doesn't have generics and since maybe this is confusing anyway, I write the logic another
	// way here.
	r.sums[n] = 0
	if r.filled > n {
		// We've already filled more than n values, so we know where to set the tail: head - n.
		r.tails[n] = (r.head + len(r.buffer) - n) % len(r.buffer)
	} else if r.filled == len(r.buffer) && n > len(r.buffer) {
		// The buffer is full and we are enlarging it. That means the tail will be the oldest entry in the
		// buffer, i.e. the one at head.
		r.tails[n] = r.head
	} else {
		// We know that n > r.filled and either the buffer isn't full or n < len(r.buffer). Thus
		// (a) The buffer isn't full. OR
		// (b) r.filled <= n <= len(r.buffer)
		// Either of these cases implies the tail is at 0, so the tail is at 0.
		r.tails[n] = 0
	}

	// Compute the sum from tail to head.
	for i := 0; i < n && i < r.filled; i++ {
		r.sums[n] += r.buffer[(r.tails[n]+i)%len(r.buffer)]
	}

	// Possibly enlarge the buffer. We assume that len(r.buffer) is being tracked because we always construct M that
	// way in New and when we start tracking timeframes in this method. Thus, we can assume here that n is strictly
	// greater than len(r.buffer).
	if n > len(r.buffer) {
		for otherN, tail := range r.tails {
			if tail >= r.head {
				r.tails[otherN] += n - len(r.buffer)
			}
		}
		padding := make([]float64, n-len(r.buffer))
		r.buffer = append(append(r.buffer[0:r.head], padding...), r.buffer[r.head:]...)
	}

	return r.sums[n] / math.Min(float64(n), float64(r.filled))
}
