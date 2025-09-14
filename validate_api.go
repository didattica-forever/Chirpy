package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	// type returnVals struct {
	// 	Valid bool `json:"valid"`
	// }	
	
	type returnCleaned struct {
		Cleaned string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}


	// validate against prophane words
	// 2. Define the list of profane words
    profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := cleanBody(params.Body, profaneWords )

	respondWithJSON(w, http.StatusOK, returnCleaned{
		Cleaned: cleanedBody,
	})
	// respondWithJSON(w, http.StatusOK, returnVals{
	// 	Valid: true,
	// })
}

// Helper function to clean the body string
func cleanBody(body string, profaneWords []string) string {
    cleaned := body
    for _, word := range profaneWords {
        // Build a regex pattern to match the word case-insensitively with word boundaries
        // \b matches word boundaries, and (?i) makes the match case-insensitive
        pattern := fmt.Sprintf(`(?i)\b%s\b`, regexp.QuoteMeta(word))
        re := regexp.MustCompile(pattern)
        
        // Replace all occurrences with "****"
        cleaned = re.ReplaceAllString(cleaned, "****")
    }
    return cleaned
}
