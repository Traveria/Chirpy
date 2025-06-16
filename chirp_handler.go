package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirps struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode paramaters", err)
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	} else {
		if len(params.Body) <= 0 {
			respondWithError(w, http.StatusBadRequest, "Empty Chirp", nil)
			return
		}
	}

	cleaned := getCleanedBody(params.Body, RestrictedWordsList)

	createParams := database.CreateChirpParams{
		Body: cleaned,
		UserID: uuid.NullUUID{
			UUID:  params.UserID,
			Valid: true,
		},
	}

	chirps, err := cfg.queries.CreateChirp(r.Context(), createParams)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error creating chirp", nil)
		return
	}

	data := Chirps{
		ID:        chirps.ID,
		CreatedAt: chirps.CreatedAt,
		UpdatedAt: chirps.UpdatedAt,
		Body:      cleaned,
		UserID:    chirps.UserID.UUID,
	}
	respondWithJson(w, http.StatusCreated, data)
}

func getCleanedBody(body string, RestrictedWordsList map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredword := strings.ToLower(word)
		if _, ok := RestrictedWordsList[loweredword]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
