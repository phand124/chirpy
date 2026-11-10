package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/phand124/chirpy/internal/database"
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
	db             *database.Queries
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

// handler to validate max length of "Chirps" and censor for language
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
		Body string `json:"cleaned_body"`
	}

	// Length check
	if len(params.Body) > 140 {
		respoondWithError(w, 400, "Chirp is too long", nil)
	} else {
		respondWithJSON(w, 200, returnVal{Body: censorPottyMouth(params.Body)})
	}

}

// filter and censor chirps
func censorPottyMouth(body string) string {
	naughtyList := [3]string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(body, " ")
	for i, word := range words {
		for _, naughty := range naughtyList {
			if strings.ToLower(word) == naughty {
				words[i] = "****"
			}
		}

	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	apiCfg := apiConfig{
		db: dbQueries,
	}

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
