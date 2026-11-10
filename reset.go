package main

import (
	"fmt"
	"net/http"
)

// Handler to reset 0
// Only accessible on "dev" platform

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respoondWithError(w, 403, "Forbidden", nil)
		return
	}
	ctx := r.Context()
	err := cfg.db.ResetUsers(ctx)
	if err != nil {
		respoondWithError(w, 500, "Error resetting users:", err)
		return
	}

	// reset Metrics for fileserverHits to 0
	message := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Swap(0))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(message))
}
