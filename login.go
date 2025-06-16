package main

import (
	"chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type logins struct {
		Password string
		Email    string
	}

	decoder := json.NewDecoder(r.Body)
	params := logins{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to decode", err)
		return
	}

	//VUD - ValidatedUserData
	VUD, err := cfg.queries.GetSingleUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(VUD.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	data := User{
		ID:        VUD.ID,
		CreatedAt: VUD.CreatedAt.Time,
		UpdatedAt: VUD.UpdatedAt.Time,
		Email:     VUD.Email,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to marshal", err)
		return
	}
	w.WriteHeader(200)
	w.Write(jsonData)
}
