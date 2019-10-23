package trading

import (
	"log"
	"time"
)

// Utility methods.
func FatalIfErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

// Window.
type Window struct {
	Begin time.Time
	End   time.Time

	Close float64
	High  float64
	Low   float64
	Open  float64
}

func (w Window) MeanPrice() float64 {
	return (w.High + w.Low) / 2
}

// Tick. Maps symbols to windows.
type Tick map[string]Window

// Position. Represents all investments and liquid assets.
type Position struct {
	Investments map[string]float64
	Liquid      float64
}

func NewPosition(liquid float64) Position {
	return Position{
		Investments: map[string]float64{},
		Liquid:      liquid,
	}

}

// Returns mature value of a position given a Tick.
func (p Position) Value(tick Tick) float64 {
	value := p.Liquid
	for symbol, quantity := range p.Investments {
		value += quantity * tick[symbol].Close
	}
	return value
}

// Strategy. For use with trading/simulate or analysis.
type Strategy interface {
	// Return slice of output headers.
	Headers() []string

	// Places initial bets.
	Initialize(tick Tick) (Position, []string)

	// Process stock tick and return position and data outputs.
	ProcessTick(tick Tick, ipos, unlist []string, returns float64) (Position, []string)
}
