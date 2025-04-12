package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg apiConfig) CreateUser(writer http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		ResponseWithError(writer, 403, "not authorized")
		return
	}

	defer req.Body.Close()

	type requestVal struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	var val requestVal
	err := decoder.Decode(&val)

	if err != nil {
		ResponseWithError(writer, 400, "error decoding json")
		return
	}
	user, err := cfg.db.CreateUser(req.Context(), sql.NullString{
		String: val.Email,
		Valid:  true,
	})

	if err != nil {
		ResponseWithError(writer, 500, err.Error())
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

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)
	writer.Write(respUser)
}
