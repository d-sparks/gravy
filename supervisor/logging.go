package supervisor

import (
	"context"
	"fmt"
	"strings"
	"time"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
	supervisor_pb "github.com/d-sparks/gravy/supervisor/proto"
	"github.com/jackc/pgx/v4"
)

var perAlgorithmCols = []string{
	"portfolio_value",
	"usd",
	"pos_value",
	"oop_value",
	"oop_deviation_min",
	"oop_deviation_max",
	"oop_deviation_10p",
	"oop_deviation_25p",
	"oop_deviation_50p",
	"oop_deviation_75p",
	"oop_deviation_90p",
	"alpha_15",
	"alpha_35",
	"alpha_252",
	"beta_15",
	"beta_35",
	"beta_252",
	"buys_value",
	"sells_value",
	"num_closing_positions",
	"num_opening_positions",
	"closing_pos_return_min",
	"closing_pos_return_max",
	"closing_pos_return_mean",
	"significant_holdings",
}

// initTimescaleDB creates a new unique identifier and creates and initializes a TimescaleDB for this session.
func (s *S) initTimescaleDB() (string, error) {
	// Connect to timescale database.
	s.timescaleContext = context.Background()
	var err error
	if s.timescaleDB, err = pgx.Connect(s.timescaleContext, s.timescaleURL); err != nil {
		return "", err
	}

	// Create identifier and table.
	s.timescaleID = fmt.Sprintf("timescaleout%d", time.Now().Unix())
	if _, err = s.timescaleDB.Exec(
		s.timescaleContext,
		fmt.Sprintf("CREATE TABLE %s (time TIMESTAMPTZ NOT NULL);", s.timescaleID),
	); err != nil {
		return "", err
	}

	return s.timescaleID, nil
}

// initTimescaleDBAlgorithmColumns, after algorithms are initialized, adds columns for each algorithm.
func (s *S) initTimescaleDBAlgorithmColumns() error {
	// Add columns for each algorithm.
	for algorithmID := range s.registrar.Algorithms {
		for _, col := range perAlgorithmCols {
			colType := "FLOAT (8)"
			if col == "significant_holdings" {
				colType = "TEXT"
			}
			algoCol := fmt.Sprintf("%s_%s", algorithmID, col)
			// Add column.
			if _, err := s.timescaleDB.Exec(
				s.timescaleContext,
				fmt.Sprintf(
					"ALTER TABLE %s ADD COLUMN %s %s NOT NULL;",
					s.timescaleID,
					algoCol,
					colType,
				),
			); err != nil {
				return err
			}
			// TODO: Potentially add indexes.
		}
	}

	return nil
}

// positionReturn calculates the return of the closing position. Also calcualtes "alternate performance" which is the
// performance of the involved assets as they were held before the position was taken. For example if the initial
// portfolio were {MSFT: 5, GOOG: 7} and one takes the position {MSFT: 4, GOOG: 20}, the initUSD will be positive to
// account for buying many more shares of GOOG (but less the price of one MSFT share which was sold). In the end, the
// performance will be based on the portfolio {MSFT: 5, GOOG: 7, initUSD: $$$} --> {MSFT: 4, GOOG: 20} while the alt
// performance is from {MSFT: 5, GOOG: 7, initUSD: $$$} --> {MSFT: 5, GOOG: 7, initUSD: $$$}. Note there is no interest
// rate included here yet.
func (s *S) positionReturn(
	portfolio *supervisor_pb.Portfolio,
	dailyData *dailyprices_pb.DailyData,
	position *Position,
) (perf float64, altPerf float64, ok bool) {
	// Calculate initial total value of assets and USD involved in the position. Also calcualte the mature value
	// had the position not been taken.
	initValue := position.initUSD
	altMatureValue := position.initUSD
	for ticker, quantity := range position.initQuantity {
		initValue += quantity * position.initPrices[ticker]
		altMatureValue += quantity * dailyData.GetPrices()[ticker].GetClose()
	}

	// Check for bogus position.
	if initValue <= 1e-5 {
		return 0.0, 0.0, false
	}

	// Calculate the mature value of the position.
	matureValue := position.initUSD
	for ticker := range position.tickers {
		matureValue += portfolio.GetStocks()[ticker] * dailyData.GetPrices()[ticker].GetClose()
	}

	// Extrapolate returns and return them.
	perf = gravy.AnnualizedPerf(initValue, matureValue, position.tradingDays)
	altPerf = gravy.AnnualizedPerf(initValue, altMatureValue, position.tradingDays)
	return perf, altPerf, true
}

// calculateClosingPositionsDistribution calculates a gravy.Distribution based on the differences in perf and altPerf
// for all closing positions.
func (s *S) calculateClosingPositionsDistribution(
	portfolio *supervisor_pb.Portfolio,
	dailyData *dailyprices_pb.DailyData,
	closingPositions map[uint64]*Position,
) *gravy.Distribution {
	values := make([]float64, len(closingPositions))
	i := 0
	for _, position := range closingPositions {
		perf, altPerf, _ := s.positionReturn(portfolio, dailyData, position)
		values[i] = perf - altPerf
	}
	lambda := func(i int) float64 { return values[i] }
	return gravy.CalculateDistribution(0, len(values), lambda)
}

// calculateOOPDistribution calculates a gravy.Distribution based on portfolio value of all stocks not in positions.
func (s *S) calculateOOPDistribution(
	portfolio *supervisor_pb.Portfolio,
	dailyData *dailyprices_pb.DailyData,
	inPositionStocks map[string]struct{},
) *gravy.Distribution {
	// Build slice of values of stocks that are not in position.
	values := []float64{}
	for ticker, quantity := range portfolio.GetStocks() {
		if _, ok := inPositionStocks[ticker]; ok {
			continue
		}
		values = append(values, quantity*dailyData.GetPrices()[ticker].GetClose())
	}

	// Calculate distribution.
	lambda := func(i int) float64 { return values[i] }
	return gravy.CalculateDistribution(0, len(values), lambda, 10, 25, 50, 75, 90)
}

// LogTick logs all data for a given tick.
func (s *S) logTickToTimescale(timestamp time.Time, dailyData *dailyprices_pb.DailyData) error {
	cols := []string{"time"}
	wildcards := []string{"$1"}
	vals := []interface{}{timestamp.Format("2006-01-02")}

	for algorithmID := range s.registrar.Algorithms {
		portfolio := s.portfolios[algorithmID]
		portfolioValue := gravy.PortfolioValue(portfolio, dailyData)
		usd := portfolio.GetUsd()

		// Calculate "in position" value
		inPositionStocks := map[string]struct{}{}
		for _, position := range s.positions[algorithmID] {
			for ticker := range position.tickers {
				inPositionStocks[ticker] = struct{}{}
			}
		}
		inPositionValue := 0.0
		for ticker := range inPositionStocks {
			inPositionValue += portfolio.GetStocks()[ticker] * dailyData.GetPrices()[ticker].GetClose()
		}

		// Calculate "out of position" value (total value less cash and in position value)
		OOPValue := portfolioValue - usd - inPositionValue
		OOPDist := s.calculateOOPDistribution(portfolio, dailyData, inPositionStocks)

		// Buys and sells value
		buysValue := s.totalBuys[algorithmID]
		s.totalBuys[algorithmID] = 0.0
		sellsValue := s.totalSells[algorithmID]
		s.totalSells[algorithmID] = 0.0

		// Closing positions
		numClosingPositions := len(s.closingPositions[algorithmID])
		closingPosDist := s.calculateClosingPositionsDistribution(
			portfolio,
			dailyData,
			s.closingPositions[algorithmID],
		)
		s.closingPositions[algorithmID] = map[uint64]*Position{}

		values := map[string]float64{
			"portfolio_value":         portfolioValue,
			"usd":                     usd,
			"pos_value":               inPositionValue,
			"oop_value":               OOPValue,
			"oop_deviation_min":       gravy.ZScore(OOPDist.Min, OOPDist.Mean, OOPDist.StDev),
			"oop_deviation_max":       gravy.ZScore(OOPDist.Max, OOPDist.Mean, OOPDist.StDev),
			"oop_deviation_10p":       gravy.ZScore(OOPDist.Percentiles[10], OOPDist.Mean, OOPDist.StDev),
			"oop_deviation_25p":       gravy.ZScore(OOPDist.Percentiles[25], OOPDist.Mean, OOPDist.StDev),
			"oop_deviation_50p":       gravy.ZScore(OOPDist.Percentiles[50], OOPDist.Mean, OOPDist.StDev),
			"oop_deviation_75p":       gravy.ZScore(OOPDist.Percentiles[75], OOPDist.Mean, OOPDist.StDev),
			"oop_deviation_90p":       gravy.ZScore(OOPDist.Percentiles[90], OOPDist.Mean, OOPDist.StDev),
			"alpha_15":                0.0, // TODO
			"alpha_35":                0.0, // TODO
			"alpha_252":               s.alpha[algorithmID].Alpha(),
			"beta_15":                 0.0, // TODO
			"beta_35":                 0.0, // TODO
			"beta_252":                s.alpha[algorithmID].Beta(),
			"buys_value":              buysValue,
			"sells_value":             sellsValue,
			"num_closing_positions":   float64(numClosingPositions),
			"num_opening_positions":   float64(s.numOpeningPositions[algorithmID]),
			"closing_pos_return_min":  closingPosDist.Min,
			"closing_pos_return_max":  closingPosDist.Max,
			"closing_pos_return_mean": closingPosDist.Mean,
		}

		for _, col := range perAlgorithmCols {
			algoCol := fmt.Sprintf("%s_%s", algorithmID, col)
			cols = append(cols, algoCol)
			switch col {
			case "significant_holdings":
				tickers, weights := gravy.SignificantAllocations(
					s.portfolios[algorithmID],
					dailyData,
					portfolioValue,
				)
				colonSeparated := make([]string, len(tickers))
				for i := 0; i < len(tickers); i++ {
					colonSeparated[i] = fmt.Sprintf("%s: %.2f", tickers[i], weights[i])
				}
				vals = append(vals, strings.Join(colonSeparated, " "))
			default:
				vals = append(vals, fmt.Sprintf("%f", values[col]))
			}

			wildcards = append(wildcards, fmt.Sprintf("$%d", len(wildcards)+1))
		}
	}

	insertQuery := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s);",
		s.timescaleID,
		strings.Join(cols, ","),
		strings.Join(wildcards, ","),
	)
	tx, err := s.timescaleDB.Begin(s.timescaleContext)
	if err != nil {
		return err
	}
	if _, err := tx.Exec(s.timescaleContext, insertQuery, vals...); err != nil {
		return err
	}

	return tx.Commit(s.timescaleContext)
}
