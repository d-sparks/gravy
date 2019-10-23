package buy_the_dip

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/d-sparks/ace-of-trades/trading"
	"github.com/d-sparks/ace-of-trades/trading/strategies/buy_and_hold"
	"github.com/gonum/stat"
)

// Strategy parameters
const sigmaBuy = 2.0
const sigmaSell = 2.0
const windowDays = 100

// BuyTheDip strategy.
type BuyTheDip struct {
	// Position.
	Position trading.Position

	// Observation is a cyclic buffer for each symbol.
	observations map[string][]float64

	// For each cyclic buffer, the index of the oldest observation.
	oldestIx map[string]int

	// Hypothetical strategies.
	BuyAndHold       buy_and_hold.BuyAndHold
	holdObservations []float64
	holdOldestIx     int
}

// Output data: number of desired buys and sells per tick.
func (b *BuyTheDip) Headers() []string {
	return []string{
		"hold_value",
		"num_buys",
		"num_sells",
		"buys",
		"sells",
	}
}

// Initializes strategy data and places initial bets.
func (b *BuyTheDip) Initialize(tick trading.Tick) (trading.Position, []string) {
	// Place initial bets.
	investment := b.Position.Liquid / float64(len(tick))
	for symbol, window := range tick {
		b.Position.Investments[symbol] = investment / window.MeanPrice()
	}
	b.Position.Liquid = 0.0

	// Initialize observations map.
	b.observations = map[string][]float64{}
	b.oldestIx = map[string]int{}
	for symbol, _ := range tick {
		b.observations[symbol] = []float64{tick[symbol].Close}
		b.oldestIx[symbol] = 0
	}

	// Initialization for hypothetical hold strategy.
	holdPosition, _ := b.BuyAndHold.Initialize(tick)
	b.holdObservations = []float64{holdPosition.Value(tick)}
	b.holdOldestIx = 0

	return b.Position, nil
}

// Hold, unless a security is no longer listed. In which case, liquidate that security.
func (b *BuyTheDip) ProcessTick(
	tick trading.Tick,
	ipos,
	unlist []string,
	returns float64,
) (trading.Position, []string) {
	// Track liquid returns.
	b.Position.Liquid += returns

	// Untrack unlisted symbols.
	for _, symbol := range unlist {
		delete(b.observations, symbol)
		delete(b.oldestIx, symbol)
		delete(b.Position.Investments, symbol)
	}

	// Prepare to observe IPOs, these will be singleton and 0 after the update below.
	for _, symbol := range ipos {
		b.observations[symbol] = []float64{}
		b.oldestIx[symbol] = 0
	}

	// Update hypothetical hold strategy.
	holdPosition, _ := b.BuyAndHold.ProcessTick(tick, ipos, unlist, returns)
	if len(b.holdObservations) < windowDays {
		b.holdObservations = append(b.holdObservations, holdPosition.Value(tick))
	} else {
		b.holdObservations[b.holdOldestIx] = holdPosition.Value(tick)
		b.holdOldestIx = (b.holdOldestIx + 1) % windowDays
	}

	// Update observations.
	for symbol, window := range tick {
		if len(b.observations[symbol]) < windowDays {
			b.observations[symbol] = append(b.observations[symbol], window.Close)
			continue
		}
		b.observations[symbol][b.oldestIx[symbol]] = window.Close
		b.oldestIx[symbol] = (b.oldestIx[symbol] + 1) % windowDays
	}

	// Slice of symbols to sell and a weighted map of dips to buy.
	totalBuy := 0.0
	potentialBuys := map[string]float64{}
	potentialSells := []string{}

	// Z-scores compared to window average for each sufficiently observed symbol.
	marketMu, marketSigma := stat.MeanStdDev(b.holdObservations, nil)
	marketZScore := stat.StdScore(holdPosition.Value(tick), marketMu, marketSigma)
	for symbol, _ := range b.Position.Investments {
		if len(b.observations[symbol]) < windowDays {
			continue
		}
		mu, sigma := stat.MeanStdDev(b.observations[symbol], nil)
		zScore := stat.StdScore(tick[symbol].Close, mu, sigma)

		marketAdjustedZScore := zScore - marketZScore

		if marketAdjustedZScore >= sigmaSell {
			potentialSells = append(potentialSells, symbol)
		} else if marketAdjustedZScore <= -sigmaBuy {
			totalBuy += marketAdjustedZScore
			potentialBuys[symbol] = marketAdjustedZScore
		}
	}

	// Logging output (number and type of buys and sells)
	dataOut := make([]string, 5)
	dataOut[0] = fmt.Sprintf("%f", b.BuyAndHold.Position.Value(tick))
	dataOut[1] = strconv.Itoa(len(potentialBuys))
	dataOut[2] = strconv.Itoa(len(potentialSells))
	if len(potentialBuys) == 0 || len(potentialSells) == 0 {
		return b.Position, dataOut
	}
	bytes, err := json.Marshal(potentialBuys)
	trading.FatalIfErr(err)
	dataOut[3] = fmt.Sprintf("%s", strings.ReplaceAll(string(bytes), ",", ";"))
	bytes, err = json.Marshal(potentialSells)
	trading.FatalIfErr(err)
	dataOut[4] = fmt.Sprintf("'%s'", strings.ReplaceAll(string(bytes), ",", ";"))

	// In reality, we'd have to leave some liquid cushion to account for fluctuations in limit
	// order pricing. For simplicity we assume we trade at exactly the most recent tick's
	// closing price.
	for _, symbol := range potentialSells {
		b.Position.Liquid += tick[symbol].Close * b.Position.Investments[symbol] * 0.25
		b.Position.Investments[symbol] *= 0.75
	}
	for symbol, weight := range potentialBuys {
		investment := b.Position.Liquid * weight / totalBuy
		b.Position.Liquid -= investment
		b.Position.Investments[symbol] += investment / tick[symbol].Close
	}

	return b.Position, dataOut
}
