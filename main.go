package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// Handler for readiness enpoint
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	message := "OK"
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(message))
}

// stateful handler for number of requests
type apiConfig struct {
	fileserverHits atomic.Int32
}

// Wrapper, increments metrics on hits to /app
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Handler for metrics, shows current count of hits
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(message))
}

// Handler to reset metrics back to 0
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Swap(0))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(message))
}

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	// main endpoint
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

	// server settings
	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serv.ListenAndServe()
}
