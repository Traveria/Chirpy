package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.queries.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retreiving chirps", nil)
		return
	}

	data := make([]Chirps, 0)

	for _, element := range chirps {
		newChrip := Chirps{
			ID:        element.ID,
			CreatedAt: element.CreatedAt,
			UpdatedAt: element.UpdatedAt,
			Body:      element.Body,
			UserID:    element.UserID.UUID,
		}
		data = append(data, newChrip)
	}

	respondWithJson(w, http.StatusOK, data)

}

func (cfg *apiConfig) handleGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	chirpID, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "chirpID invalid", nil)
		return
	}

	chirps, err := cfg.queries.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "chirp ID does not exist", nil)
		return
	}

	data := Chirps{
		ID:        chirps.ID,
		CreatedAt: chirps.CreatedAt,
		UpdatedAt: chirps.UpdatedAt,
		Body:      chirps.Body,
		UserID:    chirps.UserID.UUID,
	}

	respondWithJson(w, http.StatusOK, data)
}
