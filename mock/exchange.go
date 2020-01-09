package mock

// Exchange for simulation.
type Exchange struct {
}

// New exchange starting with a seed of USD.
func NewExchange(seed float64) *Exchange {
	return &Exchange{}
}
