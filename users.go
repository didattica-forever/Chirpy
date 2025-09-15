package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Email string `json:"email"`
}

type CreateUserResponse struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed", nil)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to read request body", err)
		return
	}

	req := CreateUserRequest{}
	if err := json.Unmarshal(body, &req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON format", err)
		return
	}

	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", err)
		return
	}

	user, err := cfg.db.CreateUser(context.Background(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusConflict, "INSERT error", err)
		return
	}

	fmt.Printf("%v", user)
	respondWithJSON(w, http.StatusCreated, CreateUserResponse{
		Id:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	})
	//respondWithJSON(w, http.StatusCreated, user)
}
