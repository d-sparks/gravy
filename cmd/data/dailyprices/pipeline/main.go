package main

import (
	"flag"
	"log"

	dailypricespipeline "github.com/d-sparks/gravy/data/dailyprices/pipeline"
)

var (
	dbURL         = flag.String("db", "postgres://localhost/gravy?sslmode=disable", "Postgres DB connection string.")
	startAtFile   = flag.String("start_at_file", "", "Skip up until the specified file.")
	startAtFolder = flag.String("start_at_folder", "nasdaq1", "Skip up until the specified folder.")
)

type folderSpec struct {
	folder   string
	exchange dailypricespipeline.Exchange
}

var folders = map[string]*folderSpec{
	"nasdaq1":      {"./data/dailyprices/raw/daily/us/nasdaq stocks/1", dailypricespipeline.NASDAQ},
	"nasdaq2":      {"./data/dailyprices/raw/daily/us/nasdaq stocks/2", dailypricespipeline.NASDAQ},
	"nasdaq_etfs":  {"./data/dailyprices/raw/daily/us/nasdaq etfs", dailypricespipeline.NASDAQ_ETF},
	"nyse1":        {"./data/dailyprices/raw/daily/us/nyse stocks/1", dailypricespipeline.NYSE},
	"nyse2":        {"./data/dailyprices/raw/daily/us/nyse stocks/2", dailypricespipeline.NYSE},
	"nyse3":        {"./data/dailyprices/raw/daily/us/nyse stocks/3", dailypricespipeline.NYSE},
	"nyse_etfs":    {"./data/dailyprices/raw/daily/us/nyse etfs", dailypricespipeline.NYSE_ETF},
	"nysemkt":      {"./data/dailyprices/raw/daily/us/nysemkt stocks", dailypricespipeline.NYSEMKT},
	"nysemkt_etfs": {"./data/dailyprices/raw/daily/us/nysemkt etfs", dailypricespipeline.NYSEMKT_ETF},
}

func main() {
	flag.Parse()

	for folderID, folder := range folders {
		if folderID < *startAtFolder {
			continue
		}
		maybeStartAtFile := ""
		if folderID == *startAtFolder {
			maybeStartAtFile = *startAtFolder
		}

		if folder.exchange == dailypricespipeline.NYSE {
			continue
		}

		log.Printf("Loading prices for exchange %s from folder `%s`\n", folder.exchange.String(), folder.folder)
		err := dailypricespipeline.Pipeline(folder.folder, *dbURL, &folder.exchange, maybeStartAtFile)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
