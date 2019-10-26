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
	dates := sort.StringSlice{}
	for date, _ := range data {
		dates = append(dates, date.Format("2006-01-02"))
	}
	dates.Sort()

	// Open output file for writing as rows of json ticks
	out, err := os.Create(output)
	gravyutil.FatalIfErr(err)
	defer out.Close()

	// Write ticks
	for _, dateStr := range dates {
		date, err := time.Parse("2006-01-02", dateStr)
		gravyutil.FatalIfErr(err)
		bytes, err := json.Marshal(data[date])
		gravyutil.FatalIfErr(err)
		tick := string(bytes) + "\n"
		_, err = out.WriteString(tick)
		gravyutil.FatalIfErr(err)
	}
}

func ParseOrDie(floatString string) float64 {
	float, err := strconv.ParseFloat(floatString, 64)
	gravyutil.FatalIfErr(err)
	return float
}
