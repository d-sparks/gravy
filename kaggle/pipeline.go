package kaggle

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/d-sparks/gravy/gravyutil"
	"github.com/d-sparks/gravy/trading"
)

func ParseOrDie(floatString string) float64 {
	float, err := strconv.ParseFloat(floatString, 64)
	gravyutil.FatalIfErr(err)
	return float
}

// Assumes the stock is published at dates[lb] and dates[ub] but not between.
// Linearly interpolate based on y = y1 + ((y2-y1)/(x2-x1))*(x-x1)
func InterpolateForSymbol(
	symbol string,
	lb int,
	ub int,
	dates []time.Time,
	data map[time.Time]*trading.Window,
) {
	// If the stock was unlisted for more than 90 days, assume the new ticker is a different
	// company.
	if ub-lb > 90 {
		return
	}

	back := data[dates[lb]]
	front := data[dates[ub]]

	dx := float64(ub - lb)
	x1 := float64(lb)

	// Slopes and y1 for each variable.
	MClose := (front.Close[symbol] - back.Close[symbol]) / dx
	MHigh := (front.High[symbol] - back.High[symbol]) / dx
	MLow := (front.Low[symbol] - back.Low[symbol]) / dx
	MOpen := (front.Open[symbol] - back.Open[symbol]) / dx
	MVolume := (front.Volume[symbol] - back.Volume[symbol]) / dx
	y1Close := back.Close[symbol]
	y1High := back.High[symbol]
	y1Low := back.Low[symbol]
	y1Open := back.Open[symbol]
	y1Volume := back.Volume[symbol]

	for i := lb + 1; i < ub; i++ {
		x := float64(i)
		window := data[dates[i]]
		window.Close[symbol] = y1Close + MClose*(x-x1)
		window.High[symbol] = y1High + MHigh*(x-x1)
		window.Low[symbol] = y1Low + MLow*(x-x1)
		window.Open[symbol] = y1Open + MOpen*(x-x1)
		window.Volume[symbol] = y1Volume + MVolume*(x-x1)
		window.Symbols.Add(symbol)
	}

}

// For days when we don't have data on a specific stock, but we do have data in the previous and
// next windows, we interpolate using the mean.
func InterpolateData(dates []time.Time, data map[time.Time]*trading.Window) {
	// Get the last index of when a symbol was observed.
	lastListingIx := map[string]int{}
	for i := 0; i < len(dates); i++ {
		for symbol, _ := range data[dates[i]].Symbols {
			lastListingIx[symbol] = i
		}
	}

	for i := 1; i+1 < len(dates); i++ {
		wPrev := data[dates[i-1]]
		w := data[dates[i]]
		missing := wPrev.Symbols.Minus(w.Symbols)

		for symbol, _ := range missing {
			// No need to do anything for unlisted symbols.
			if i >= lastListingIx[symbol] {
				continue
			}

			// Find the next index for which this symbol is listed.
			j := i + 1
			for ; !data[dates[j]].Symbols.Contains(symbol); j++ {
			}

			// So we have listings at index i-1 and j, and we need to fill [i, j-1].
			InterpolateForSymbol(symbol, i-1, j, dates, data)
		}
	}
}

func Pipeline(input, output string) {
	// Open input file for reading as CSV
	f, err := os.Open(input)
	gravyutil.FatalIfErr(err)
	reader := csv.NewReader(f)

	// Populate headers
	headers, err := reader.Read()
	gravyutil.FatalIfErr(err)
	index := map[string]int{}
	for i, header := range headers {
		index[header] = i
	}
	get := func(row []string, key string) string {
		return row[index[key]]
	}

	// Group by date
	data := map[time.Time]*trading.Window{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		gravyutil.FatalIfErr(err)

		dateStr := get(row, "date")
		date, err := time.Parse("2006-01-02", dateStr)
		gravyutil.FatalIfErr(err)
		if _, ok := data[date]; !ok {
			data[date] = &trading.Window{
				Open:   trading.Prices{},
				Close:  trading.Prices{},
				High:   trading.Prices{},
				Low:    trading.Prices{},
				Volume: trading.Prices{},
			}
		}
		ticker := get(row, "ticker")
		data[date].Begin = date
		data[date].End = date
		data[date].Open[ticker] = ParseOrDie(get(row, "open"))
		data[date].Close[ticker] = ParseOrDie(get(row, "close"))
		data[date].High[ticker] = ParseOrDie(get(row, "high"))
		data[date].Low[ticker] = ParseOrDie(get(row, "low"))
		data[date].Volume[ticker] = ParseOrDie(get(row, "volume"))
	}

	// Get tick order
	dateStrings := sort.StringSlice{}
	for date, _ := range data {
		dateStrings = append(dateStrings, date.Format("2006-01-02"))
	}
	dateStrings.Sort()
	dates := make([]time.Time, len(dateStrings))
	for i, dateString := range dateStrings {
		date, err := time.Parse("2006-01-02", dateString)
		gravyutil.FatalIfErr(err)
		dates[i] = date
	}

	// Interpolate data for missing symbols.
	InterpolateData(dates, data)

	// Open output file for writing as rows of json ticks
	out, err := os.Create(output)
	gravyutil.FatalIfErr(err)
	defer out.Close()

	// Write ticks
	for _, date := range dates {
		bytes, err := json.Marshal(data[date])
		gravyutil.FatalIfErr(err)
		tick := string(bytes) + "\n"
		_, err = out.WriteString(tick)
		gravyutil.FatalIfErr(err)
	}
}
