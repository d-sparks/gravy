package db

import (
	"time"

	"github.com/d-sparks/gravy/trading"
)

type Data struct {
	Window trading.Window
}

type Store interface {
	Get(date time.Time) Data
}
