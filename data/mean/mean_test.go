package mean

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMean tests the mean compared to a non-streaming implementation.
func TestMean(t *testing.T) {
	xs := []float64{3.0, 4.0, 5.0, 17.0, -1.0, 600.0}
	sum := 0.0
	mu := NewStreaming()

	for _, x := range xs {
		sum += x
		mu.Observe(x)
	}

	assert.Equal(t, sum/float64(len(xs)), mu.Value())
}
