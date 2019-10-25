package main

import (
	"bufio"
	"flag"
	"log"
	"strconv"

	"github.com/d-sparks/gravy/algorithm"
	"github.com/d-sparks/gravy/db/dailywindow"
	"github.com/d-sparks/gravy/exchange"
)

var windows = flag.String("windows", "./data/kaggle/historical_as_ticks.json", "Kaggledata")
var symbols = flag.String("symbols", "./data/kaggle/historical_stocks.csv", "Stock symbols")
var output = flag.String("output", "./results", "Results output")

// Default stores, typically in memory stores.
func GetDataStores(dailywindowFilename string) map[string]Store {
	stores := map[string]Store{}
	stores[dailywindow.Name] = dailywindow.NewInMemoryStore(dailywindowFilename)
	return stores
}

// Writes a CSV line from a key order and kv. Includes an integer ID line.
func WriteCSVLine(id int, order []string, kv map[string]string, out bufio.Writer) {
	out.WriteString(strconv.Itoa(id))
	for _, header := range order {
		out.WriteString("," + kv[header])
	}
	out.WriteString("\n")
}

// Simulates over all dates from a dailwindow.InMemoryStore.
func Simulate(stores map[string]Store, seed float64, output string) {
	// Get dates to simulate.
	dailywindow, ok := stores[dailywindow.Name].(dailywindow.InMemoryStore)
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

	// Iterate over dates and simulate trading.
	const hideAfterIndex int = len(dates) / 2
	for i, date := range dates {
		algorithm.Trade(date)
		const hide bool = i > hideAfterIndex
		WriteCsvLine(i, algorithm.Headers(), algorithm.Debug(hide), out)
	}
}

func main() {
	flag.Parse()
	stores := GetDataStores(*windows)
	Simulate(stores, 1.0, *output)
}
