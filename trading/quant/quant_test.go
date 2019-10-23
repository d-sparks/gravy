package quant

import (
	"math/rand"
	"testing"

	"github.com/gonum/stat"
	"github.com/stretchr/testify/assert"
)

const epsilon float64 = 1E-10

func TestStatStream_MuSigma(t *testing.T) {
	statStream := NewStatStream()
	rand := rand.New(rand.NewSource(1337))

	// Make a stream of random numbers, observe them, and calculate the typical mean.
	data := make([]float64, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = rand.Float64()
		statStream.Observe(data[i])
	}

	mu, sigma := stat.MeanStdDev(data, nil)
	assert.InDelta(t, mu, statStream.Mu(), epsilon)
	assert.InDelta(t, sigma, statStream.Sigma(), epsilon)
}
