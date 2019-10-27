package main

import (
	"bufio"
	"log"
	"strconv"

	"github.com/d-sparks/gravy/algorithm"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/gravyutil"
	"github.com/d-sparks/gravy/mock"
	"github.com/spf13/cobra"
)

var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate a trading session",
	Run:   simulateFn,
}

var windows string
var symbols string
var output string

func init() {
	rootCmd.AddCommand(simulateCmd)
	simulateCmd.Flags().StringVarP(&windows, "windows", "w", "./data/kaggle/historical_as_windows.json", "Kaggledata")
	simulateCmd.Flags().StringVarP(&symbols, "symbols", "s", "./data/kaggle/historical_stocks.csv", "Stock symbols")
	simulateCmd.Flags().StringVarP(&output, "output", "o", "./results", "Results output")
}

func simulateFn(cmd *cobra.Command, args []string) {
	stores := GetDataStores(windows)
	Simulate(stores, 1.0, output)
}

// Default stores, typically in memory stores.
func GetDataStores(dailywindowFilename string) map[string]db.Store {
	stores := map[string]db.Store{}
	stores[dailywindow.Name] = dailywindow.NewInMemoryStore(dailywindowFilename)
	return stores
}

// Writes a CSV header.
func WriteCSVHeader(headers []string, out *bufio.Writer) {
	out.WriteString("id")
	for _, header := range headers {
		out.WriteString("," + header)
	}
	out.WriteString("\n")
}

// Writes a CSV line from a key order and kv. Includes an integer ID line.
func WriteCSVLine(id int, order []string, kv map[string]string, out *bufio.Writer) {
	out.WriteString(strconv.Itoa(id))
	for _, header := range order {
		out.WriteString("," + kv[header])
	}
	out.WriteString("\n")
}

// Simulates over all dates from a dailwindow.InMemoryStore.
func Simulate(stores map[string]db.Store, seed float64, output string) {
	// Get dates to simulate.
	dailywindow, ok := stores[dailywindow.Name].(*dailywindow.InMemoryStore)
	if !ok {
		log.Fatalf("Simulation failed, dailywindow store not an InMemoryStore")
	}
	dates := dailywindow.Dates()

	// Mock exchange.
	exchange := mock.NewExchange(seed)

	// Make trading algorithm.
	algorithm := algorithm.NewTradingAlgorithm(stores, exchange)

	// Create output file.
	out := gravyutil.FileWriterOrDie(output)

	// Iterate over dates and simulate trading, export CSV.
	WriteCSVHeader(algorithm.Headers(), out)
	skipUntilIndex := 3650
	hideAfterIndex := len(dates) / 2
	for i := skipUntilIndex; i < len(dates); i++ {
		algorithm.Trade(dates[i])
		hide := i > hideAfterIndex
		WriteCSVLine(i, algorithm.Headers(), algorithm.Debug(hide), out)
	}
}
