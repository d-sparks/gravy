package pairstatspipeline

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"github.com/d-sparks/gravy/data/covariance"
	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/data/mean"
	"github.com/d-sparks/gravy/data/variance"
	"github.com/d-sparks/gravy/registrar"
	"github.com/golang/protobuf/ptypes"
	timestamp_pb "github.com/golang/protobuf/ptypes/timestamp"
)

type Mode int

const (
	Count Mode = iota
)

func fatalIfErr(err error) {
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func parseTimeOrDie(timeString string) *timestamp_pb.Timestamp {
	nativeTime, err := time.Parse("2006-01-02", timeString)
	fatalIfErr(err)
	timestamp, err := ptypes.TimestampProto(nativeTime)
	fatalIfErr(err)
	return timestamp
}

// CartesianChunk represents an indexed portion of a partesian product: [xMin, xMax) x [yMin, yMax).
type CartesianChunk struct {
	xMin int
	xMax int
	yMin int
	yMax int
}

// NewCartesianChunk makes a "square" cartesian chunk: equal number of indexes on x and y.
func NewCartesianChunk(xMin, yMin, chunkSize int) *CartesianChunk {
	return &CartesianChunk{xMin: xMin, xMax: xMin + chunkSize, yMin: yMin, yMax: yMin + chunkSize}
}

// String returns a `[a, b) x [c, d)` string representation.
func (c *CartesianChunk) String() string {
	return fmt.Sprintf("[%d, %d) x [%d, %d)", c.xMin, c.xMax, c.yMin, c.yMax)
}

// chunkResult is the output of processing a CartesianChunk.
type ChunkResult struct {
	count int
}

// ProcessChunk does the thing.
func ProcessChunk(
	ctx context.Context,
	r *registrar.R,
	c *CartesianChunk,
	tradingDates *dailyprices_pb.TradingDates,
	assetIDs []int64,
	idToTickerExchangePairs map[int64]*dailyprices_pb.AssetIdsResponse_TickerExchangePair,
	mode Mode,
) (*ChunkResult, error) {
	const period int = 252

	// Make rolling covariance trackers for each pair, as well as variance of correlation tracker. We create
	// separate mean trackers for each pair in case some assets pricings are sparse - we think of this as separate
	// random variables for each pair of assets but only defined when both prices are available.
	xMu := map[int64]map[int64]*mean.Rolling{}
	yMu := map[int64]map[int64]*mean.Rolling{}
	xVar := map[int64]map[int64]*variance.Rolling{}
	yVar := map[int64]map[int64]*variance.Rolling{}
	cov := map[int64]map[int64]*covariance.Rolling{}
	corMu := map[int64]map[int64]*mean.Rolling{}
	corVar := map[int64]map[int64]*variance.Rolling{}
	for x := c.xMin; x < c.xMax; x++ {
		xID := assetIDs[x]
		xMu[xID] = map[int64]*mean.Rolling{}
		yMu[xID] = map[int64]*mean.Rolling{}
		xVar[xID] = map[int64]*variance.Rolling{}
		yVar[xID] = map[int64]*variance.Rolling{}
		cov[xID] = map[int64]*covariance.Rolling{}
		corMu[xID] = map[int64]*mean.Rolling{}
		corVar[xID] = map[int64]*variance.Rolling{}
		for y := c.yMin; y < c.yMax; y++ {
			yID := assetIDs[y]
			xMu[xID][yID] = mean.NewRolling(period)
			yMu[xID][yID] = mean.NewRolling(period)
			xVar[xID][yID] = variance.NewRolling(xMu[xID][yID], period)
			yVar[xID][yID] = variance.NewRolling(yMu[xID][yID], period)
			cov[xID][yID] = covariance.NewRolling(xMu[xID][yID], yMu[xID][yID], period)
			corMu[xID][yID] = mean.NewRolling(period)
			corVar[xID][yID] = variance.NewRolling(corMu[xID][yID], period)
		}
	}

	// Iterate through trading dates
	chunkResult := ChunkResult{count: 0}
	for ix, tradingDate := range tradingDates.GetTimestamps() {
		// Get prices
		dailyData, err := r.DailyPrices.Get(ctx, &dailyprices_pb.Request{Timestamp: tradingDate, Version: 0})
		if err != nil {
			return nil, fmt.Errorf("Error getting prices: %s", err.Error())
		}

		// Maybe log progress
		if (c.xMin == 0 && c.yMin == 0 && ix%252 == 0) || ix%(5*252) == 0 {
			log.Printf("...processing date %s...\n", tradingDate.AsTime().Format("2006-01-02"))
		}

		// Iterate through pairs
		for xID := range cov {
			// Get x prices or skip if not listed.
			xTicker := idToTickerExchangePairs[xID].GetTicker()
			xPrices, ok := dailyData.GetPrices()[xTicker]
			if !ok {
				continue
			}
			xPrice := xPrices.GetClose()

			for yID := range cov[xID] {
				// Get y prices or skip if not listed.
				yTicker := idToTickerExchangePairs[yID].GetTicker()
				yPrices, ok := dailyData.GetPrices()[yTicker]
				if !ok {
					continue
				}
				yPrice := yPrices.GetClose()

				// Record the statistics.
				xMu[xID][yID].Observe(xPrice)
				yMu[xID][yID].Observe(yPrice)
				xVar[xID][yID].Observe(xPrice)
				yVar[xID][yID].Observe(yPrice)
				cov[xID][yID].Observe(xPrice, yPrice)
				correlation :=
					cov[xID][yID].Value() /
						math.Sqrt(xVar[xID][yID].Value()*yVar[xID][yID].Value())
				corMu[xID][yID].Observe(correlation)
				corVar[xID][yID].Observe(correlation)

				// If an inclusion criterion is met, do the thing
				if !cov[xID][yID].Full() {
					continue
				}

				chunkResult.count += 1
			}
		}
	}

	return &chunkResult, nil
}

func CalculateCorrelations(mode Mode) error {
	// Create a registrar.
	log.Println("Creating registrar...")
	r, err := registrar.New()
	if err != nil {
		return err
	}

	// Initiate a new session.
	log.Println("Starting new session...")
	ctx, cancel := context.WithTimeout(context.Background(), 2.0*time.Hour)
	defer cancel()
	req := dailyprices_pb.NewSessionRequest{
		SimRange: &dailyprices_pb.Range{
			Lb: parseTimeOrDie("1900-01-01"),
			Ub: parseTimeOrDie("2100-01-01"),
		},
	}
	_, err = r.DailyPrices.NewSession(ctx, &req)
	if err != nil {
		return err
	}

	// Get asset ID map.
	log.Println("Getting asset ids...")
	assetIDs, err := r.DailyPrices.AssetIds(ctx, &dailyprices_pb.AssetIdsRequest{})
	if err != nil {
		return err
	}
	if len(assetIDs.GetAssetIds()) == 0 {
		return fmt.Errorf("No asset ids.")
	}

	// Get trading dates in range.
	log.Println("Getting all trading dates...")
	tradingDates, err := r.DailyPrices.TradingDatesInRange(ctx, req.GetSimRange())
	if err != nil {
		return err
	}
	if len(tradingDates.GetTimestamps()) == 0 {
		return fmt.Errorf("No trading dates.")
	}

	// Get sorted list of ids.
	log.Println("Sorting ids...")
	ids := make([]int64, len(assetIDs.GetAssetIds()))
	idIx := 0
	for id := range assetIDs.GetAssetIds() {
		ids[idIx] = id
		idIx++
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	// Chunk ids. Targeting 10GB, and using twice the estimate of a covariance size:
	//
	//   936 ~= sqrt(10*1024*1024*1024 / (2*( 8(3*252 + 3 + 1) + 4(7 + 3))))
	//
	// so we would want to take 936 x 936 chunks, or 876096 pairs at a time. However, this does end up using enough
	// memory to swap on 32GB of RAM (with other processes running; and the above estimate doesn't quite include
	// everything). So we decrease this by 50% by multiplying the chunk size by \sqrt(0.5), giving 661, or 436921
	// pairs at a time.
	const chunkSize int = 661

	// Inner loop body executes a total of 256 times per 10000 assets.
	count := 0
	for xMin := 0; xMin < len(ids); xMin += chunkSize {
		for yMin := 0; yMin < len(ids); yMin += chunkSize {
			chunk := NewCartesianChunk(xMin, yMin, chunkSize)
			log.Printf("Processing chunk %s...\n", chunk.String())
			chunkResult, err := ProcessChunk(ctx, r, chunk, tradingDates, ids, assetIDs.GetAssetIds(), mode)
			if err != nil {
				return err
			}
			count += chunkResult.count
			log.Printf("... count so far: %d\n", count)
		}
	}

	log.Printf("\n==================================\nOutput:\n\n")
	log.Printf("Count: %d\n", count)
	log.Println("\n==================================")

	return nil
}
