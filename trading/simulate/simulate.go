package simulate

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Clever/go-utils/stringset"
	"github.com/d-sparks/ace-of-trades/trading"
)

// Write a string to a file or die
func writeOrDie(f *os.File, s string) {
	_, err := f.WriteString(s)
	trading.FatalIfErr(err)
}

// Headers: date, value, liquid, [strategy columns], [position value columns]
func headerString(dataHeaders, symbols []string) string {
	dataHeadersString := strings.Join(dataHeaders, ",")
	symbolHeadersString := strings.Join(symbols, ",")
	return fmt.Sprintf("date,value,liquid,%s,%s\n", dataHeadersString, symbolHeadersString)
}

// Log a CSV row corresponding to headers above.
func logString(tick trading.Tick, position trading.Position, data, symbols []string) string {
	// Hacky date
	var date string
	for _, v := range tick {
		date = v.Begin.Format("2006-01-02")
		break
	}

	// Value and data columns
	value := position.Value(tick)
	liquid := position.Liquid
	dataColumns := strings.Join(data, ",")

	// Position value columns
	values := make([]string, len(symbols))
	for i, symbol := range symbols {
		value := tick[symbol].Close * position.Investments[symbol]
		values[i] = fmt.Sprintf("%f", value)
	}
	valueColumns := strings.Join(values, ",")

	return fmt.Sprintf("%s,%f,%f,%s,%s\n", date, value, liquid, dataColumns, valueColumns)
}

// Returns keys in tick that are not in other.
func keyDiff(tick, other trading.Tick) []string {
	diff := []string{}
	for symbol, _ := range tick {
		if _, ok := other[symbol]; !ok {
			diff = append(diff, symbol)
		}
	}
	return diff
}

// In the event of unlisting, assume payout of position * last close price.
func unlistingReturns(symbols []string, position trading.Position, price trading.Tick) float64 {
	returns := 0.0
	for _, symbol := range symbols {
		returns += price[symbol].Close * position.Investments[symbol]
	}
	return returns
}

// Creates a bufio.Scanner with large buffer given a filename.
func newScanner(filename string) *bufio.Scanner {
	file, err := os.Open(filename)
	trading.FatalIfErr(err)
	scanner := bufio.NewScanner(file)
	scanner.Buffer([]byte{}, 1024*1024)
	return scanner
}

// Parses next line of a bufio.Scanner and unmarshals into a Tick.
func parseNextTick(scanner *bufio.Scanner) trading.Tick {
	tick := trading.Tick{}
	trading.FatalIfErr(json.Unmarshal(scanner.Bytes(), &tick))
	return tick
}

// Reads the tick file and gathers all symbols ever mentioned.
func symbolsFromFile(filename string) []string {
	scanner := newScanner(filename)
	symbols := stringset.New()
	scanner.Scan() // Scan past header row.
	for scanner.Scan() {
		line := string(scanner.Bytes())
		columns := strings.Split(line, ",")
		if len(columns) == 0 {
			log.Fatalf("Bad read in symbols CSV line %s", line)
		}
		symbols.Add(columns[0])
	}
	trading.FatalIfErr(scanner.Err())
	return symbols.ToList()
}

// Replays a newline separated, JSON encoded trading.Ticks through a trading.Strategy.
func SimulateFromFile(
	ticksFilename,
	symbolsFilename,
	outputFilename string,
	strategy trading.Strategy,
) {
	// Setup file i/o.
	scanner := newScanner(ticksFilename)
	out, err := os.Create(outputFilename)
	trading.FatalIfErr(err)
	defer out.Close()

	// Prepare to scan input file, track previous tick and position, read symbols.
	var previousTick trading.Tick
	var previousPosition trading.Position
	var data []string

	// Read symbols from scanner and reset the file.
	symbols := symbolsFromFile(symbolsFilename)

	// Write headers.
	writeOrDie(out, headerString(strategy.Headers(), symbols))

	// Skip ten years worth of data due to sparsity.
	for i := 0; i < 3650; i++ {
		scanner.Scan()
	}

	// Initialize from first scan.
	if scanner.Scan() {
		previousTick = parseNextTick(scanner)
		previousPosition, data = strategy.Initialize(previousTick)
		writeOrDie(out, logString(previousTick, previousPosition, data, symbols))
	}

	// Replay ticks and write output data.
	for scanner.Scan() {
		tick := parseNextTick(scanner)
		ipo := keyDiff(tick, previousTick)
		unlist := keyDiff(previousTick, tick)
		returns := unlistingReturns(unlist, previousPosition, previousTick)

		position, data := strategy.ProcessTick(tick, ipo, unlist, returns)
		writeOrDie(out, logString(tick, position, data, symbols))

		previousTick = tick
	}
	trading.FatalIfErr(scanner.Err())
}
