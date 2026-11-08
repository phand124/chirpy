package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	message := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html")
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

// handler to validate max length of "Chirps"
func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respoondWithError(w, 500, "Error decoding parameters:", err)
		return
	}

	type returnVal struct {
		Valid bool `json:"valid"`
	}

	if len(params.Body) > 140 {
		respoondWithError(w, 400, "Chirp is too long", nil)
	} else {
		respondWithJSON(w, 200, returnVal{Valid: true})
	}

}

// generic response with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

// gerneic response with error
func respoondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Printf("Error: %s", err)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	errMsg := errorResponse{Error: msg}

	respondWithJSON(w, code, errMsg)
}

func main() {
	mux := http.NewServeMux()
	apiCfg := apiConfig{}

	// main endpoint
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// server settings
	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	serv.ListenAndServe()
}
