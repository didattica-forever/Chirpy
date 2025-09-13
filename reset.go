package main

import "net/http"

// to display the number of hits
func (cfg *apiConfig) resetStatsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// str := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	// w.Write([]byte(str))
	w.Write([]byte("Hits reset to 0"))
}
