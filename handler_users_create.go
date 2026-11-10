package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respoondWithError(w, 500, "Error decoding parameters:", err)
		return
	}

	user, err := cfg.db.CreateUser(ctx, params.Email)
	if err != nil {
		respoondWithError(w, 500, "Error creating user:", err)
		return
	}

	respondWithJSON(w, 201, User(user))

}
