package main

import (
	"flag"
)

var ticks = flag.String("ticks", "./data/kaggle/historical_as_ticks.json", "Kaggledata")
var symbols = flag.String("symbols", "./data/kaggle/historical_stocks.csv", "Stock symbols")
var output = flag.String("output", "./results", "Results output")

//func main() {
//	strategy := buy_the_dip.BuyTheDip{
//		Position: trading.NewPosition(1.0),
//		BuyAndHold: buy_and_hold.BuyAndHold{
//			Position: trading.NewPosition(1.0),
//		},
//	}
//	simulate.SimulateFromFile(*ticks, *symbols, *output, &strategy)
//}
