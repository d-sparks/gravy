package main

import (
	"bufio"
	"flag"
	"log"
	"strconv"

	"github.com/d-sparks/gravy/algorithm"
	"github.com/d-sparks/gravy/db"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/exchange"
	"github.com/d-sparks/gravy/gravyutil"
)

var windows = flag.String("windows", "./data/kaggle/historical_as_windows.json", "Kaggledata")
var symbols = flag.String("symbols", "./data/kaggle/historical_stocks.csv", "Stock symbols")
var output = flag.String("output", "./results", "Results output")

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
	exchange := exchange.NewMockExchange(seed)

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

func main() {
	flag.Parse()
	stores := GetDataStores(*windows)
	Simulate(stores, 1.0, *output)
}
