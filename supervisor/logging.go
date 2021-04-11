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
			algoCol := fmt.Sprintf("%s_%s", algorithmID, col)
			// Add column.
			if _, err := s.timescaleDB.Exec(
				s.timescaleContext,
				fmt.Sprintf(
					"ALTER TABLE %s ADD COLUMN %s FLOAT (8) NOT NULL;",
					s.timescaleID,
					algoCol,
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

	ix := 2
	for algorithmID := range s.registrar.Algorithms {
		for _, col := range perAlgorithmCols {
			algoCol := fmt.Sprintf("%s_%s", algorithmID, col)
			cols = append(cols, algoCol)
			switch col {
			case "usd":
				vals = append(vals, fmt.Sprintf("%f", s.portfolios[algorithmID].GetUsd()))
			case "portfolio_value":
				vals = append(
					vals,
					fmt.Sprintf("%f", gravy.PortfolioValue(s.portfolios[algorithmID], dailyData)),
				)
			case "alpha_252":
				vals = append(vals, fmt.Sprintf("%f", s.alpha[algorithmID].Alpha()))
			case "beta_252":
				vals = append(vals, fmt.Sprintf("%f", s.alpha[algorithmID].Beta()))
			default:
				vals = append(vals, "0.0")
			}

			wildcards = append(wildcards, fmt.Sprintf("$%d", ix))
			ix += 1
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
