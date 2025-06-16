package main

import (
	"chirpy/internal/auth"
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handleNewUser(w http.ResponseWriter, r *http.Request) {
	type NewUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := NewUser{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters", err)
		return
	}

	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed,
	}

	user, err := cfg.queries.CreateUser(r.Context(), userParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "sql error", err)
		return
	}

	data := User{ID: user.ID, CreatedAt: user.CreatedAt.Time, UpdatedAt: user.UpdatedAt.Time, Email: user.Email}
	jsonData, err := json.Marshal(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters", err)
		return
	}

	w.WriteHeader(201)
	w.Write(jsonData)

}
