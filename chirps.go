package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sambakker4/chirpy/internal/auth"
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
	token, err := auth.GetBearerToken(req.Header)

	if err != nil {
		ResponseWithError(writer, 500, "Token not recieved")
	}

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)

	if err != nil {
		ResponseWithError(writer, 500, "Error decoding json")
		return
	}

	id, err := auth.ValidateJWT(token, cfg.tokenSecret)

	if err != nil {
		ResponseWithError(writer, 401, "Error validating token")
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
		UserID: id,
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

func (cfg apiConfig) GetAllChirps(writer http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(req.Context())
	writer.Header().Set("Content-Type", "application/json")
	if err != nil {
		ResponseWithError(writer, 500, "Error retrieving all chirps from database")
		return
	}

	chirps := make([]Chirp, 0)

	for _, chirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	ResponseWithJson(writer, 200, chirps)
}

func (cfg apiConfig) GetChirp(writer http.ResponseWriter, req *http.Request) {
	chirpIDstring := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstring)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		ResponseWithError(writer, 500, "Error parsing chirp id")
		return
	}

	dbChirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		ResponseWithError(writer, 404, "chirp not found")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	ResponseWithJson(writer, 200, chirp)
}

func (cfg apiConfig) DeleteChirp(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	chirpIDstring := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstring)

	if err != nil {
		ResponseWithError(writer, 500, "Error parsing chirp id")
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ResponseWithError(writer, 401, "Error retrieving token from header")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		ResponseWithError(writer, 500, "Error validating access token")
		return
	}

	chirp, err := cfg.db.GetChirp(req.Context(), chirpID)
	if err != nil {
		ResponseWithError(writer, 404, "Chirp not found")
		return
	}

	if chirp.UserID != userID {
		ResponseWithError(writer, 403, "Not authorized")
		return
	}
	
	err = cfg.db.DeleteChirp(req.Context(), chirpID)
	if err != nil {
		ResponseWithError(writer, 400, "Error deleting chirp")
		return
	}
	
	writer.WriteHeader(204)
	writer.Write([]byte(""))
}
