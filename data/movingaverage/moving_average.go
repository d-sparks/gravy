package movingaverage

import (
	"math"
)

// M is a struct that tracks n-day averages for an asset or portfolio.
type M struct {
	buffer []float64
	sums   map[int]float64
	head   int
	tails  map[int]int
	filled int
}

// New creates a moving average struct that tracks the given windows.
func New(days ...int) *M {
	m := M{sums: map[int]float64{}, head: 0, tails: map[int]int{}, filled: 0}
	maxRequest := 0
	for _, n := range days {
		if n <= 0 {
			panic("Can't track 0 day average.")
		}
		if n > maxRequest {
			maxRequest = n
		}
		m.sums[n] = 0.0
		m.tails[n] = 0
	}
	m.buffer = make([]float64, maxRequest)
	return &m
}

// Observe inserts the value. The way this works is there is one large cyclic buffer whose size is equal to the longest
// moving average being tracked, however there are several tails where the shorter moving averages end.
func (m *M) Observe(value float64) {
	// Update the sums and tails if necessary.
	for n := range m.sums {
		if m.filled >= n {
			m.sums[n] -= m.buffer[m.tails[n]]
			m.tails[n] = (m.tails[n] + 1) % len(m.buffer)
		}
		m.sums[n] += value
	}

	// Insert into cyclic buffer and advance the head index.
	m.buffer[m.head] = value
	m.head = (m.head + 1) % len(m.buffer)
	if m.filled < len(m.buffer) {
		m.filled++
	}
}

// Value returns the moving average for the prescribed timeframe. Begins tracking the timeframe if it is not already
// being tracked. Assumes n is > 0 and that there has been at least one observation.
func (m *M) Value(n int) float64 {
	if n <= 0 {
		panic("Can't track 0 day average.")
	}

	// Return the average if it is already tracked.
	if sum, ok := m.sums[n]; ok {
		return sum / math.Min(float64(n), float64(m.filled))
	}

	// Otherwise, begin tracking. The tail will be at m.head - n % len(m.buffer) unless the buffer is not yet big
	// enough or not yet full enough. When the buffer is not full enough, e.g. n >= m.filled, we want the tail to be
	// at 0, or m.head - m.filled. When the buffer is not big enough, we want to start at the oldest element, which
	// is at m.head, i.e. we want m.head - len(m.buffer) = m.head. This can all be expressed as
	//
	//   m.tails[n] = (m.head + len(m.buffer) - min(n, m.filled, len(m.buffer)) % len(m.buffer)
	//
	// However, since Go doesn't have generics and since maybe this is confusing anyway, I write the logic another
	// way here.
	m.sums[n] = 0
	if m.filled > n {
		// We've already filled more than n values, so we know where to set the tail: head - n.
		m.tails[n] = (m.head + len(m.buffer) - n) % len(m.buffer)
	} else if m.filled == len(m.buffer) && n > len(m.buffer) {
		// The buffer is full and we are enlarging it. That means the tail will be the oldest entry in the
		// buffer, i.e. the one at head.
		m.tails[n] = m.head
	} else {
		// We know that n > m.filled and either the buffer isn't full or n < len(m.buffer). Thus
		// (a) The buffer isn't full. OR
		// (b) m.filled <= n <= len(m.buffer)
		// Either of these cases implies the tail is at 0, so the tail is at 0.
		m.tails[n] = 0
	}

	// Compute the sum from tail to head.
	for i := 0; i < n && i < m.filled; i++ {
		m.sums[n] += m.buffer[(m.tails[n]+i)%len(m.buffer)]
	}

	// Possibly enlarge the buffer. We assume that len(m.buffer) is being tracked because we always construct M that
	// way in New and when we start tracking timeframes in this method. Thus, we can assume here that n is strictly
	// greater than len(m.buffer).
	if n > len(m.buffer) {
		for otherN, tail := range m.tails {
			if tail >= m.head {
				m.tails[otherN] += n - len(m.buffer)
			}
		}
		padding := make([]float64, n-len(m.buffer))
		m.buffer = append(append(m.buffer[0:m.head], padding...), m.buffer[m.head:]...)
	}

	return m.sums[n] / math.Min(float64(n), float64(m.filled))
}
