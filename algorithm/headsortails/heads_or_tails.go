package headsortails

import (
	"bufio"
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"

	algorithmio_pb "github.com/d-sparks/gravy/algorithm/proto"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
	"github.com/d-sparks/gravy/registrar"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/golang/protobuf/ptypes"
)

// Mode of operation: training or inference
type Mode int

const (
	// Training mode generates training data.
	Training Mode = iota

	// Inference mode actually trades.
	Inference
)

// Prediction is the result: heads or tails.
type Prediction int

const (
	// Heads is the favorable prediction.
	Heads = iota

	// Tails is the unfavorable prediction.
	Tails
)

// String is for writing to CSV.
func (p Prediction) String() string {
	if p == Heads {
		return "1.0"
	}
	return "0.0"
}

// MinResults is the minimum number of observed results before generating training data or trading.
const MinResults int = 15

// HeadsOrTails is an algorithm that will try to predict whether each asset will go up or down and invests accordingly.
type HeadsOrTails struct {
	algorithmio_pb.UnimplementedAlgorithmServer

	// Algorithm ID (usually "buyandhold" unless multiple are running)
	id          string
	algorithmID *supervisor_pb.AlgorithmId
	registrar   *registrar.R

	// Non data state
	mode          Mode
	samplingRatio float64
	numResults    int
	rand          *rand.Rand
	sink          *bufio.Writer

	// Data
	recentResults        map[string][]Prediction
	previousTickPrices   map[string]*dailyprices_pb.Prices
	previousTickFeatures map[string]*features
}

// skipTrading is a precondition. Save time if you don't need to fetch prices/portfolio.
func (b *HeadsOrTails) skipTrading() bool {
	return false
}

// headsCondition is the ground truth condition that we try to predict.
func headsCondition(prev, curr *dailyprices_pb.Prices) Prediction {
	if curr.GetClose() > prev.GetClose() && curr.GetClose() > curr.GetOpen() {
		return Heads
	}
	return Tails
}

// features of the heads/tails model. These are meant to be all relative, thus comparable for different stocks at
// different times.
type features struct {
	recentMovement []Prediction
	zVol15         float64
	z15            float64
	z35            float64
	z252           float64
	beta           float64
	sigmaMarket15  float64
	sigmaMarket35  float64
	sigmaMarket252 float64
	zMarket15      float64
	zMarket35      float64
	zMarket252     float64
}

// extractFeatures extracts features from a dailyprices_pb.DailyData.
func extractFeatures(ticker string, data *dailyprices_pb.DailyData) *features {
	// Make sure we have this asset's and market's prices & stats.
	prices, pricesOK := data.GetPrices()[ticker]
	stats, statsOK := data.GetStats()[ticker]
	marketPrices, marketPricesOK := data.GetPrices()["SPY"]
	marketStats, marketStatsOK := data.GetStats()["SPY"]
	if !(pricesOK && statsOK && marketPricesOK && marketStatsOK) {
		return nil
	}

	// Things that we'll use more than twice.
	price := prices.GetClose()
	marketPrice := marketPrices.GetClose()
	sigmaMarket15 := math.Sqrt(marketStats.GetMovingVariance()[15])
	sigmaMarket35 := math.Sqrt(marketStats.GetMovingVariance()[35])
	sigmaMarket252 := math.Sqrt(marketStats.GetMovingVariance()[252])

	// Return features.
	return &features{
		zVol15: gravy.ZScore(
			prices.GetVolume(),
			stats.GetMovingVolume()[15],
			math.Sqrt(stats.GetMovingVolumeVariance()[15]),
		),
		z15: gravy.ZScore(
			price,
			stats.GetMovingAverages()[15],
			math.Sqrt(stats.GetMovingVariance()[15]),
		),
		z35: gravy.ZScore(
			price,
			stats.GetMovingAverages()[35],
			math.Sqrt(stats.GetMovingVariance()[35]),
		),
		z252: gravy.ZScore(
			price,
			stats.GetMovingAverages()[252],
			math.Sqrt(stats.GetMovingVariance()[252]),
		),
		beta:           stats.GetBeta(),
		sigmaMarket15:  sigmaMarket15,
		sigmaMarket35:  sigmaMarket35,
		sigmaMarket252: sigmaMarket252,
		zMarket15: gravy.ZScore(
			marketPrice,
			marketStats.GetMovingAverages()[15],
			sigmaMarket15,
		),
		zMarket35: gravy.ZScore(
			marketPrice,
			marketStats.GetMovingAverages()[35],
			sigmaMarket35,
		),
		zMarket252: gravy.ZScore(
			marketPrice,
			marketStats.GetMovingAverages()[252],
			sigmaMarket252,
		),
	}
}

// headers of the output CSV. (There's one more header: the label, as last column.)
var headers []string = []string{
	"recent",
	"zVol15",
	"z15",
	"z35",
	"z252",
	"beta",
	"sigmaMarket15",
	"sigmaMarket35",
	"sigmaMarket252",
	"zMarket15",
	"zMarket35",
	"zMarket252",
}

// getters for cells of output CSV.
var getters = map[string]func(f *features) string{
	"recent":         func(f *features) string { return "0.0" },
	"zVol15":         func(f *features) string { return fmt.Sprintf("%f", f.zVol15) },
	"z15":            func(f *features) string { return fmt.Sprintf("%f", f.z15) },
	"z35":            func(f *features) string { return fmt.Sprintf("%f", f.z35) },
	"z252":           func(f *features) string { return fmt.Sprintf("%f", f.z252) },
	"beta":           func(f *features) string { return fmt.Sprintf("%f", f.beta) },
	"sigmaMarket15":  func(f *features) string { return fmt.Sprintf("%f", f.sigmaMarket15) },
	"sigmaMarket35":  func(f *features) string { return fmt.Sprintf("%f", f.sigmaMarket35) },
	"sigmaMarket252": func(f *features) string { return fmt.Sprintf("%f", f.sigmaMarket252) },
	"zMarket15":      func(f *features) string { return fmt.Sprintf("%f", f.zMarket15) },
	"zMarket35":      func(f *features) string { return fmt.Sprintf("%f", f.zMarket35) },
	"zMarket252":     func(f *features) string { return fmt.Sprintf("%f", f.zMarket252) },
}

func (b *HeadsOrTails) emitExamples(data *dailyprices_pb.DailyData) error {
	if b.previousTickFeatures == nil {
		// If this is the first tick, write headers.
		_, err := b.sink.WriteString(strings.Join(append(headers, "result"), ",") + "\n")
		if err != nil {
			return fmt.Errorf("Error writing headers: %s", err.Error())
		}
		b.recentResults = map[string][]Prediction{}
		b.previousTickPrices = map[string]*dailyprices_pb.Prices{}
		b.previousTickFeatures = map[string]*features{}
	} else if b.numResults >= MinResults {
		// If sufficient observations have been made, emit training data.
		for ticker, currPrices := range data.GetPrices() {
			// Subsample.
			if b.rand.Float64() >= b.samplingRatio {
				continue
			}

			// Only record example if we have successive ticks of data.
			prevPrices, pricesOK := b.previousTickPrices[ticker]
			prevFeatures, featuresOK := b.previousTickFeatures[ticker]
			if !(pricesOK && featuresOK && prevPrices != nil && prevFeatures != nil) {
				continue
			}

			// Make columns from previous tick features.
			cols := make([]string, len(headers)+1)
			for i, header := range headers {
				cols[i] = getters[header](prevFeatures)
			}

			// Record result.
			cols[len(headers)] = headsCondition(prevPrices, currPrices).String()

			// Write example.
			if _, err := b.sink.WriteString(strings.Join(cols, ",") + "\n"); err != nil {
				return fmt.Errorf("Error writing observation %d: %s", b.numResults, err.Error())
			}
		}
	}

	// Track data.
	b.numResults++
	for ticker, prices := range data.GetPrices() {
		b.previousTickPrices[ticker] = prices
		b.previousTickFeatures[ticker] = extractFeatures(ticker, data)
	}

	return b.sink.Flush()
}

// trade is the algorithm itself.
func (b *HeadsOrTails) trade(
	portfolio *supervisor_pb.Portfolio,
	data *dailyprices_pb.DailyData,
) ([]*supervisor_pb.Order, error) {
	if b.mode == Training {
		return nil, b.emitExamples(data)
	}
	// Note: Inference is not implemented in Golang. Use the python algorithm (./heads_or_tails.py) for inference.
	return nil, nil
}

// New creates a new, uninitialized HeadsOrTails algorithm.
func New(algorithmID string) *HeadsOrTails {
	return &HeadsOrTails{
		id:          algorithmID,
		algorithmID: &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
		mode:        Inference,
	}
}

// NewForTraining creates a new, uninitialized HeadsOrTails algorithm.
func NewForTraining(algorithmID string, samplingRatio float64, sampleSeed int64, sink *bufio.Writer) *HeadsOrTails {
	return &HeadsOrTails{
		id:            algorithmID,
		algorithmID:   &supervisor_pb.AlgorithmId{AlgorithmId: algorithmID},
		mode:          Training,
		samplingRatio: samplingRatio,
		rand:          rand.New(rand.NewSource(sampleSeed)),
		sink:          sink,
	}
}

// ******************************
//  Mostly boilerplate hereafter
// ******************************

// Init initializes the registrar. The algorithm should be listening before calling Init to avoid deadlocks.
func (b *HeadsOrTails) Init() error {
	var err error
	b.registrar, err = registrar.NewWithSupervisor()
	return err
}

// Close closes the regitsrar.
func (b *HeadsOrTails) Close() {
	b.registrar.Close()
}

// Execute implements the algorithm interface.
func (b *HeadsOrTails) Execute(ctx context.Context, input *algorithmio_pb.Input) (*algorithmio_pb.Output, error) {
	fmt.Printf("Excuting algorithm on %s\n", ptypes.TimestampString(input.GetTimestamp()))
	orders := []*supervisor_pb.Order{}

	if !b.skipTrading() {
		portfolio, err := b.registrar.Supervisor.GetPortfolio(ctx, b.algorithmID)
		if err != nil {
			return nil, fmt.Errorf("Error getting portfolio in `%s`: %s", b.id, err.Error())
		}

		req := dailyprices_pb.Request{Timestamp: input.GetTimestamp(), Version: 0}
		prices, err := b.registrar.DailyPrices.Get(ctx, &req)
		if err != nil {
			return nil, fmt.Errorf("Error getting daily prices in `%s`: %s", b.id, err.Error())
		}

		if orders, err = b.trade(portfolio, prices); err != nil {
			return nil, fmt.Errorf("Error trading or writing training data: %s", err.Error())
		}
		for _, order := range orders {
			if _, err := b.registrar.Supervisor.PlaceOrder(ctx, order); err != nil {
				return nil, fmt.Errorf(
					"Error placing order from `%s`: %s", b.id, err.Error(),
				)
			}
		}
	}
	if _, err := b.registrar.Supervisor.DoneTrading(ctx, b.algorithmID); err != nil {
		return nil, fmt.Errorf("Error calling DoneTrading from `%s`: %s", b.id, err.Error())
	}

	return &algorithmio_pb.Output{}, nil
}
