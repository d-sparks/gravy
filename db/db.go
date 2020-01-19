package db

import (
	"time"

	"github.com/Clever/go-utils/stringset"
)

type Data struct {
	TickersToPrices map[string]Prices
	Tickers         stringset.StringSet
}

type Store interface {
	ValidDate(date time.Time) (bool, error)
	Get(date time.Time) (*Data, error)
	NextDate(date time.Time) (*time.Time, error)
}

type Prices struct {
	Open     float64
	Close    float64
	AdjClose float64
	Low      float64
	High     float64
	Volume   float64
}
