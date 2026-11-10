package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
