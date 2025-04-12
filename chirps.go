package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sambakker4/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) CreateChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		ResponseWithError(writer, 500, "Error decoding json")
		return

	}

	message := params.Body

	if len(message) > 140 {
		ResponseWithError(writer, 400, "Chirp is too long")
		return
	}
	message = removeBadWords(message)

	dbChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
		Body:   message,
		UserID: params.UserID,
	})

	if err != nil {
		ResponseWithError(writer, 500, "Error retrieving data from database")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	ResponseWithJson(writer, 201, chirp)
}

func removeBadWords(str string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	words := strings.Split(str, " ")

	for i, word := range words {
		if idx := slices.Index(badWords, strings.ToLower(word)); idx != -1 {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
