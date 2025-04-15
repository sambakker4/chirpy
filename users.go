package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sambakker4/chirpy/internal/auth"
	"github.com/sambakker4/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg apiConfig) CreateUser(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	if cfg.platform != "dev" {
		ResponseWithError(writer, 403, "not authorized")
		return
	}

	defer req.Body.Close()

	type requestVal struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	var val requestVal
	err := decoder.Decode(&val)

	if err != nil {
		ResponseWithError(writer, 400, "error decoding json")
		return
	}

	hashed_password, err := auth.HashPassword(val.Password)
	if err != nil {
		ResponseWithError(writer, 500, "error hashing password")
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email: sql.NullString{
			String: val.Email,
			Valid:  true,
		},
		HashedPassword: hashed_password,
	})

	if err != nil {
		ResponseWithError(writer, 500, "Error creating user")
		return
	}

	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email.String,
	}

	respUser, err := json.Marshal(newUser)
	if err != nil {
		ResponseWithError(writer, 500, "")
	}

	writer.WriteHeader(201)
	writer.Write(respUser)
}
