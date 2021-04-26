package main

import (
	"log"

	pairstatspipeline "github.com/d-sparks/gravy/data/pairstats/pipeline"
)

func main() {
	if err := pairstatspipeline.CalculateCorrelations(pairstatspipeline.Count); err != nil {
		log.Fatalf(err.Error())
	}
}
