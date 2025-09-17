package main

import (
	"Chirpy/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"Chirpy/internal/auth"
)

type ChirpMsg struct {
	Body string    `json:"body"`
	Id   uuid.UUID `json:"user_id"`
}

type CreateChirpResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) chirpGetHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	// id := r.URL.Query().Get("chirpID") non funziona perchè non è una variabile

	fmt.Printf("==>> id: [%s]\n", id)

	if id == "" {
		respondWithError(w, http.StatusBadRequest, "ID is required in the URL query, e.g., /api/chirps?id=myid", nil)
		return
	}

	uId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error parsing id", err)
		return
	}

	chirp, err := cfg.db.GetChirp(context.Background(), uId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp Not Found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, CreateChirpResponse{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *apiConfig) chirpListHandler(w http.ResponseWriter, r *http.Request) {
	chirpDbList, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusConflict, "INSERT error", err)
		return
	}

	chirpList := []CreateChirpResponse{}
	for _, chirp := range chirpDbList {
		chirpList = append(chirpList, CreateChirpResponse{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpList)
}

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing Bearer Token", err)
		return
	}

	userid, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Bearer Token", err)
		return
	}

	fmt.Printf("==>> userid: [%s]\n", userid)

	decoder := json.NewDecoder(r.Body)

	chirp := ChirpMsg{}
	err = decoder.Decode(&chirp)
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
	cleanedBody := cleanBody(chirp.Body, profaneWords)

	parms := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userid, //chirp.Id,
	}

	chirpDb, err := cfg.db.CreateChirp(context.Background(), parms)
	if err != nil {
		respondWithError(w, http.StatusConflict, "INSERT error", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, CreateChirpResponse{
		Id:        chirpDb.ID,
		CreatedAt: chirpDb.CreatedAt,
		UpdatedAt: chirpDb.UpdatedAt,
		Body:      chirpDb.Body,
		UserID:    chirpDb.UserID,
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
