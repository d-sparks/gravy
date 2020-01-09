package signal

import (
	"fmt"
	"time"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/gravy/db"
)

// Data output by signals.
type SignalOutput struct {
	KV        map[string]float64
	StringSet stringset.StringSet
}

// Signals compute data to be used to inform a Strategy or TradingAlgorithm.
type Signal interface {
	// Return the identifier for this signal.
	Name() string

	// Compute and/or return cached signal output.
	Compute(date time.Time, stores map[string]db.Store) (*SignalOutput, error)

	// Debug info for previous computation.
	Headers() []string
	Debug() map[string]string
}

// Most signals will want to wrap themselves in a CachedSignal for simplicity. This way, a signal
// can assume its Compute method will only be called once per timestamp.
type CachedSignal struct {
	cache       map[time.Time]*SignalOutput
	signal      Signal
	expireAfter time.Duration
}

func NewCachedSignal(signal Signal, expireAfter time.Duration) *CachedSignal {
	return &CachedSignal{
		cache:       map[time.Time]*SignalOutput{},
		signal:      signal,
		expireAfter: expireAfter,
	}
}

// Evict oldest cache entries.
func (c *CachedSignal) Evict(date time.Time) {
	for cacheDate, _ := range c.cache {
		if date.Sub(cacheDate).Nanoseconds() > c.expireAfter.Nanoseconds() {
			delete(c.cache, cacheDate)
		}
	}
}

func (c *CachedSignal) Name() string {
	return c.signal.Name()
}

// Compute and add to cache or serve from cache. Also evicts the cache.
func (c *CachedSignal) Compute(date time.Time, stores map[string]db.Store) (*SignalOutput, error) {
	c.Evict(date)
	if _, ok := c.cache[date]; !ok {
		signalOutput, err := c.signal.Compute(date, stores)
		if err != nil {
			return nil, fmt.Errorf("Error computing signal `%s`: `%s`", c.signal.Name(), err.Error())
		}
		c.cache[date] = signalOutput
	}
	return c.cache[date], nil
}

// Forwarding methods.
func (c *CachedSignal) Headers() []string        { return c.signal.Headers() }
func (c *CachedSignal) Debug() map[string]string { return c.signal.Debug() }
