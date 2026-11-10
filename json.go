package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// generic response with JSON
func respondWithJSON(w http.ResponseWriter, code int, payload any) {

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
