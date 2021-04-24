package supervisor

import (
	"context"
	"fmt"
	"strings"
	"time"

	dailyprices_pb "github.com/d-sparks/gravy/data/dailyprices/proto"
	"github.com/d-sparks/gravy/gravy"
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
		oopValue := portfolioValue - usd - inPositionValue

		// Buys and sells value
		buysValue := s.totalBuys[algorithmID]
		s.totalBuys[algorithmID] = 0.0
		sellsValue := s.totalSells[algorithmID]
		s.totalSells[algorithmID] = 0.0

		// Closing positions
		numClosingPositions := len(s.closingPositions[algorithmID])
		s.closingPositions[algorithmID] = map[uint64]*Position{}

		values := map[string]float64{
			"portfolio_value":         portfolioValue,
			"usd":                     usd,
			"pos_value":               inPositionValue,
			"oop_value":               oopValue,
			"oop_deviation_min":       0.0, // TODO
			"oop_deviation_max":       0.0, // TODO
			"oop_deviation_10p":       0.0, // TODO
			"oop_deviation_25p":       0.0, // TODO
			"oop_deviation_50p":       0.0, // TODO
			"oop_deviation_75p":       0.0, // TODO
			"oop_deviation_90p":       0.0, // TODO
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
			"closing_pos_return_min":  0.0, // TODO
			"closing_pos_return_max":  0.0, // TODO
			"closing_pos_return_mean": 0.0, // TODO
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
