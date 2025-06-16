package main

import (
	"net/http"
)

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.PLATFORM == "dev" {
		_, err := cfg.queries.DeleteAllUsers(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "sql delete error", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)

	} else {
		w.WriteHeader(http.StatusForbidden)
		return
	}
}
