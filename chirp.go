package main

import (
	"Chirpy/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"regexp"

	"github.com/google/uuid"
	//"Chirpy/internal/database"
)

type ChirpMsg struct {
  Body string `json:"body"`
  Id uuid.UUID `json:"user_id"`
}

type CreateChirpResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserID     uuid.UUID    `json:"user_id"`
}

func (cfg *apiConfig)  chirpHandler(w http.ResponseWriter, r *http.Request) {

	type returnCleaned struct {
		Cleaned string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)

	chirp := ChirpMsg{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode Chirp", err)
		return
	}

	const maxChirpLength = 140
	if len(chirp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}


	// validate against prophane words
	// 2. Define the list of profane words
    profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanedBody := cleanBody(chirp.Body, profaneWords )

	// respondWithJSON(w, http.StatusOK, returnCleaned{
	// 	Cleaned: cleanedBody,
	// })
	// // respondWithJSON(w, http.StatusOK, returnVals{
	// // 	Valid: true,
	// // })

	parms := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: chirp.Id,
	}

	chirpDb, err := cfg.db.CreateChirp(context.Background(), parms)
	if err != nil {
		respondWithError(w, http.StatusConflict, "INSERT error", err)
		return
	}

	fmt.Printf("%v\n", chirpDb)
	respondWithJSON(w, http.StatusCreated, CreateChirpResponse{
		Id:        chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt.Time,
		UpdatedAt: chirpDb.UpdatedAt.Time,
		Body:     chirpDb.Body,
		UserID:     chirpDb.UserID,
	})
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
