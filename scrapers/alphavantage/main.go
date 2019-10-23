package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

// Based on Alphavantage API spec and guidelines (5 requests per minute)
const intraday = "TIME_SERIES_INTRADAY"
const interval5min = "5min"
const requestInterval = 12 * time.Second

type alphavantageRequest struct {
	hostname string
	function string
	symbol   string
	interval string
	apikey   string
}

// Date formatting
const fmtYYYYMMDD = "2006-01-02"

// Program flags
var hostname = flag.String("hostname", "http://localhost:8080", "alphavantage hostname")
var apikey = flag.String("apikey", "bogus", "alphavantage apikey")
var symbolsFile = flag.String("symbols", "data/sp500", "filepath to symbols to scrape")
var outputDir = flag.String("outputdir", "data/alphavantage", "filepath for output")

// Returns body of the HTTP response or an error.
func retrieve(avReq alphavantageRequest) (io.ReadCloser, error) {
	queryParams := url.Values{}
	queryParams.Add("function", avReq.function)
	queryParams.Add("symbol", avReq.symbol)
	queryParams.Add("interval", avReq.interval)
	queryParams.Add("apikey", avReq.apikey)

	req := fmt.Sprintf("%s?%s", avReq.hostname, queryParams.Encode())
	log.Printf("Sending request `%s`...\n", req)

	resp, err := http.Get(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP response `%d`", resp.StatusCode)
	}
	return resp.Body, nil
}

// Writes the contents of given readcloser to file. Closes readcloser if file write was attempted.
func writeToFile(filename string, body io.ReadCloser) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		return fmt.Errorf("Cannot write `%s`, file already exists", filename)
	}
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	defer body.Close()
	_, err = io.Copy(file, body)
	return err
}

// Loads the target symbols from a newline separated file or returns an error.
func loadSymbolsFile(filename string, symbols *[]string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		*symbols = append(*symbols, scanner.Text())
	}

	return scanner.Err()
}

// Scrapes for the set of symbols, putting results in outputDir/YYYY-MM-DD/*. Can retry n times.
// Uses the existence of files to track whether that symbol's scrape is done.
func scrape(date string, symbols []string, interval time.Duration, retries int) error {
	log.Printf("\n\nBegin scraping %s\n\n", date)
	incomplete := false

	// Make target date directory if it doesn't exist
	err := os.MkdirAll(path.Join(*outputDir, date), 0700)
	if err != nil {
		return err
	}

	for _, symbol := range symbols {
		log.Printf("Scraping %s...\n", symbol)

		// Skip if we've already scraped this symbol
		filename := path.Join(*outputDir, date, symbol)
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			log.Printf("Already scraped %s on %s", symbol, date)
			continue
		}

		// Hit API
		avReq := alphavantageRequest{
			hostname: *hostname,
			function: intraday,
			symbol:   symbol,
			interval: interval5min,
			apikey:   *apikey,
		}
		respBody, err := retrieve(avReq)
		time.Sleep(interval)
		if err != nil {
			log.Printf("... %s request failed: `%s`\n", symbol, err.Error())
			incomplete = true
			continue
		}

		// Write file
		if err = writeToFile(filename, respBody); err != nil {
			log.Printf("... %s file write failed: `%s`\n", symbol, err.Error())
			incomplete = true
		}
	}

	// Retry logic
	if incomplete {
		if retries > 0 {
			log.Print("Scrape failed, retrying...\n")
			return scrape(date, symbols, interval, retries-1)
		}
		return fmt.Errorf("Scrape failed")
	}

	return nil
}

func main() {
	flag.Parse()

	// Load symbols
	symbols := []string{}
	if err := loadSymbolsFile(*symbolsFile, &symbols); err != nil {
		log.Fatalf("Error loading symbols: %s", err.Error())
	}

	// Scrape symbols (no retries)
	date := time.Now().Format(fmtYYYYMMDD)
	if err := scrape(date, symbols, requestInterval, 0); err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("Scrape successful")
}
