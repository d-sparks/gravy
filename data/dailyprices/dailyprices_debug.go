package dailyprices

import (
	"fmt"
	"net/http"
	// _ "net/http/pprof"	// For debug
)

func (s *Server) RunDebugServer(port int) {
	http.HandleFunc("/", s.HandleDebug)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func humanReadable(KBSize int) string {
	if KBSize > 1024 {
		MBSize := KBSize / 1024
		if MBSize > 1024 {
			return fmt.Sprintf("%d GB", MBSize/1024)
		}
		return fmt.Sprintf("%d MB", MBSize)
	}
	return fmt.Sprintf("%d KB", KBSize)
}

func (s *Server) HandleDebug(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.times) == 0 {
		fmt.Fprintf(w, "\n\n\tNo data")
		return
	}

	fmt.Fprintf(
		w,
		"\n\n\tTimes:\t\t%s - %s",
		s.times[0].Format("2006-01-02"),
		s.times[len(s.times)-1].Format("2006-01-02"),
	)
	fmt.Fprintf(w, "\n\n\tCache entries:\t%d", len(s.cache[0]))
	numAssets := len(s.stats)
	fmt.Fprintf(w, "\n\n\tStats for:\t%d assets", numAssets)
	const approximateKBPerAsset int = 55
	fmt.Fprintf(w, "\n\tAprox mem:\t%s", humanReadable(numAssets*approximateKBPerAsset))

	numPairStats := 0
	for _, second := range s.pairStats {
		numPairStats += len(second)
	}
	fmt.Fprintf(w, "\n\n\tPairStats:\t%d pairs of assets", numPairStats)
	const approximateBytesPerPair int = 144
	fmt.Fprintf(w, "\n\tApprox mem:\t%s", humanReadable(numPairStats*approximateBytesPerPair/1024))
}

/*
	How to estimate KB per Stat.

	mean size = (
		x floats
		n floats
		1 int
		n ints
		1 int
	) = 8(x + n) + 4(2 + n)

	covariance size = 3 * mean.size  + float + int
			= 8(3x + 3n + 1) + 4(7 + 3n)

	alpha = 2 * covariance size + 3 float + int
	      = 8(6x + 6n + 5) + 4(15 + 6n)



	stats = alpha(252, 3) + 3 * mean(252, 3) + 2 * (covariance(252, 1) + covariance(35, 1) + covariance(15, 1))
	        + 20*8

		= 8(6*252 + 6*3 + 5) + 4(15 + 6*3) + 3(8(3*252 + 3*3 + 1) + 4(7 + 3*3)) +
		  2 * (8(3*252 + 3*1 + 1) + 4(7 + 3) + 8(3*35 + 3*1 + 1) + 4(7 + 3) + 8(3*15 + 3*1 + 1) + 4(7 + 3)) +
		  20*8
		= 46076
*/
