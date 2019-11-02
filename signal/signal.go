package signal

import (
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
	// Compute and/or return cached signal output.
	Compute(date time.Time, stores map[string]db.Store) SignalOutput

	// Debug info for previous computation.
	Headers() []string
	Debug() map[string]string
}

// Most signals will want to wrap themselves in a CachedSignal for simplicity. This way, a signal
// can assume its Compute method will only be called once per timestamp.
type CachedSignal struct {
	cache       map[time.Time]SignalOutput
	signal      Signal
	expireAfter time.Duration
}

func NewCachedSignal(signal Signal, expireAfter time.Duration) *CachedSignal {
	return &CachedSignal{
		cache:       map[time.Time]SignalOutput{},
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

// Compute and add to cache or serve from cache. Also evicts the cache.
func (c *CachedSignal) Compute(date time.Time, stores map[string]db.Store) SignalOutput {
	c.Evict(date)
	if _, ok := c.cache[date]; !ok {
		c.cache[date] = c.signal.Compute(date, stores)
	}
	return c.cache[date]
}

// Forwarding methods.
func (c *CachedSignal) Headers() []string        { return c.signal.Headers() }
func (c *CachedSignal) Debug() map[string]string { return c.signal.Debug() }
