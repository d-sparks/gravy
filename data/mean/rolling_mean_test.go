package mean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMovingAverageThreeAverages tests a moving average covering 1, 2, and 5 day spacing including overflow.
func TestMovingAverageThreeAverages(t *testing.T) {
	m := NewRolling(1, 2, 5)

	// Insert a first element and check identical averages.
	m.Observe(17.0)
	assert.Equal(t, 17.0, m.Value(1))
	assert.Equal(t, 17.0, m.Value(2))
	assert.Equal(t, 17.0, m.Value(5))

	// Insert a second element: now the 1 day average should diverge.
	m.Observe(13.0)
	assert.Equal(t, 13.0, m.Value(1))
	assert.Equal(t, 15.0, m.Value(2))
	assert.Equal(t, 15.0, m.Value(5))

	// Insert three more elements to fill the buffer.
	for i := 2; i < 5; i++ {
		m.Observe(float64(i))
	}
	assert.Equal(t, 4.0, m.Value(1))
	assert.Equal(t, (3.0+4.0)/2.0, m.Value(2))
	assert.Equal(t, (17.0+13.0+2.0+3.0+4.0)/5.0, m.Value(5))

	// Now overflow the buffer and check that things work as expected.
	for i := 500; i < 517; i++ {
		m.Observe(float64(i))
	}
	assert.Equal(t, 516.0, m.Value(1))
	assert.Equal(t, (515.0+516.0)/2.0, m.Value(2))
	assert.Equal(t, (512.0+513.0+514.0+515.0+516.0)/5.0, m.Value(5))
}

// TestMovingAverageTrackNewRollingValues tests the case of tracking a new timeframe by asking for a value that is not yet
// tracked.
func TestMovingAverageTrackNewRollingValues(t *testing.T) {
	m := NewRolling(3)

	// Insert numerous values and ask for a smaller and much larger timeframe.
	for i := 50; i < 55; i++ {
		m.Observe(float64(i))
	}
	assert.Equal(t, (53.0+54.0)/2.0, m.Value(2))
	assert.Equal(t, (52.0+53.0+54.0)/3.0, m.Value(3))
	assert.Equal(t, (52.0+53.0+54.0)/3.0, m.Value(7))

	// Fill up to exactly 7 cyclic values and check.
	for i := 0; i < 4; i++ {
		m.Observe(float64(i))
	}
	assert.Equal(t, 3.0, m.Value(1))
	assert.Equal(t, (2.0+3.0)/2.0, m.Value(2))
	assert.Equal(t, (1.0+2.0+3.0)/3, m.Value(3))
	assert.Equal(t, (53.0+54.0+1.0+2.0+3.0)/6.0, m.Value(6))
	assert.Equal(t, (52.0+53.0+54.0+0.0+1.0+2.0+3.0)/7.0, m.Value(7))

	// Insert many more elements and see that things are as expected.
	for i := 10; i < 20; i++ {
		m.Observe(float64(i))
	}
	assert.Equal(t, (18.0+19.0)/2.0, m.Value(2))
	assert.Equal(t, (14.0+15.0+16.0+17.0+18.0+19.0)/6.0, m.Value(6))
	assert.Equal(t, (13.0+14.0+15.0+16.0+17.0+18.0+19.0)/7.0, m.Value(7))
	assert.Equal(t, (13.0+14.0+15.0+16.0+17.0+18.0+19.0)/7.0, m.Value(100))
}

// TestOldestValue checks that the oldest value helper works.
func TestOldestValue(t *testing.T) {
	m := NewRolling(5)

	// Insert a first element and check identical averages.
	m.Observe(17.0)
	assert.Equal(t, 17.0, m.OldestValue())

	// Check again after many observations.
	for i := 500; i < 517; i++ {
		m.Observe(float64(i))
	}
	assert.Equal(t, 512.0, m.OldestValue())
}
