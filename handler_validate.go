package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChiprsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnvals struct {
		CleanedBody string `json:"cleaned_body"`
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
	}

	cleaned := getCleanedBody(params.Body, RestrictedWordsList)

	respondWithJson(w, http.StatusOK, returnvals{
		CleanedBody: cleaned,
	})
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
