package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/d-sparks/gravy/algorithm"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailyprices"
	"github.com/d-sparks/gravy/mock"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate a trading session",
	Run:   simulateFn,
}

var (
	dbURL       string
	pricesTable string
	begin       string
	end         string
	output      string
)

func init() {
	rootCmd.AddCommand(simulateCmd)
	f := simulateCmd.Flags()
	f.StringVarP(&dbURL, "dburl", "d", "postgres://localhost/gravy?sslmode=disable", "Location of Postgres gravydb")
	f.StringVarP(&pricesTable, "pricestable", "t", "dailyprices", "Table with dailyprices")
	f.StringVarP(&begin, "begin", "b", "1337-01-23", "Begin date.")
	f.StringVarP(&end, "end", "e", "2337-01-23", "End date.")
	f.StringVarP(&output, "output", "o", "./gravy.csv", "Output csv file.")
}

func simulateFn(cmd *cobra.Command, args []string) {
	// Connect to Postgres.
	dailypricesStore, err := dailyprices.NewPostgresStore(dbURL, pricesTable)
	if err != nil {
		log.Fatalf("Couldn't create stores: `%s`", err.Error())
	}
	defer dailypricesStore.Close()

	// Build stores.
	stores := map[string]db.Store{}
	stores[dailyprices.Name] = dailypricesStore

	// Get trading dates.
	allDates, err := dailypricesStore.AllDates()
	if err != nil {
		log.Fatalf("Error getting trading dates: `%s`", err.Error())
	}

	if err = Simulate(allDates, stores, 1.0, output); err != nil {
		log.Fatalf("Error during simulation: `%s`", err.Error())
	}
}

// Writes a CSV header.
func WriteCSVHeader(headers []string, out *os.File) {
	out.WriteString("id")
	for _, header := range headers {
		out.WriteString("," + header)
	}
	out.WriteString("\n")
}

// Writes a CSV line from a key order and kv. Includes an integer ID line.
func WriteCSVLine(id int, order []string, kv map[string]string, out *os.File) {
	out.WriteString(strconv.Itoa(id))
	for _, header := range order {
		out.WriteString("," + kv[header])
	}
	out.WriteString("\n")
}

// Simulates over specified dates and runs the trading algorithm, directing debug header output to the given writer.
func Simulate(allDates []time.Time, stores map[string]db.Store, seed float64, output string) error {
	// Mock exchange.
	exchange := mock.NewExchange(seed)

	// Make trading algorithm.
	algorithm := algorithm.NewTradingAlgorithm(stores, exchange)

	// Create output file.
	out, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("Error opening output file: `%s`", err.Error())
	}
	defer out.Close()

	// Iterate over dates and simulate trading, export CSV.
	WriteCSVHeader(algorithm.Headers(), out)
	for i := 0; i < len(allDates); i++ {
		dateStr := allDates[i].Format("2006-01-02")
		if dateStr > end {
			break
		}
		if dateStr < begin {
			continue
		}
		if err := algorithm.Trade(allDates[i]); err != nil {
			return fmt.Errorf("Error in trading algorithm: `%s`", err.Error())
		}
		WriteCSVLine(i, algorithm.Headers(), algorithm.Debug( /*hide=*/ false), out)
	}
	return nil
}
