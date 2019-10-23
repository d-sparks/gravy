package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/d-sparks/ace-of-trades/trading"
)

var input = flag.String("input", "./data/kaggle/historical_stock_prices.csv", "Kaggledata input")
var output = flag.String("output", "./data/kaggle/historical_as_ticks.json", "Normalized output")

func ParseOrDie(floatString string) float64 {
	float, err := strconv.ParseFloat(floatString, 64)
	trading.FatalIfErr(err)
	return float
}

func main() {
	// Open input file for reading as CSV
	f, err := os.Open(*input)
	trading.FatalIfErr(err)
	reader := csv.NewReader(f)

	// Populate headers
	headers, err := reader.Read()
	trading.FatalIfErr(err)
	index := map[string]int{}
	for i, header := range headers {
		index[header] = i
	}
	get := func(row []string, key string) string {
		return row[index[key]]
	}

	// Group by date
	data := map[string]trading.Tick{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		trading.FatalIfErr(err)

		date := get(row, "date")
		time, err := time.Parse("2006-01-02", date)
		trading.FatalIfErr(err)
		if _, ok := data[date]; !ok {
			data[date] = trading.Tick{}
		}
		data[date][get(row, "ticker")] = trading.Window{
			Close: ParseOrDie(get(row, "close")),
			High:  ParseOrDie(get(row, "high")),
			Low:   ParseOrDie(get(row, "low")),
			Open:  ParseOrDie(get(row, "open")),
			Begin: time,
			End:   time,
		}
	}

	// Get tick order
	dates := sort.StringSlice{}
	for date, _ := range data {
		dates = append(dates, date)
	}
	dates.Sort()

	// Open output file for writing as rows of json ticks
	out, err := os.Create(*output)
	trading.FatalIfErr(err)
	defer out.Close()

	// Write ticks
	for _, date := range dates {
		bytes, err := json.Marshal(data[date])
		trading.FatalIfErr(err)
		tick := string(bytes) + "\n"
		_, err = out.WriteString(tick)
		trading.FatalIfErr(err)
	}
}
