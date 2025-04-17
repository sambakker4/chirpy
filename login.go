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

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg apiConfig) Login(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	writer.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	var loginInfo LoginInfo
	err := decoder.Decode(&loginInfo)

	if err != nil {
		ResponseWithError(writer, 400, "Error decoding json")
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), sql.NullString{
		Valid:  true,
		String: loginInfo.Email,
	})

	if err != nil {
		ResponseWithError(writer, 500, "Error retrieving user from database")
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, loginInfo.Password)
	if err != nil {
		ResponseWithError(writer, 401, "incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(
		user.ID, cfg.tokenSecret, time.Second*time.Duration(60*60),
	)

	if err != nil {
		ResponseWithError(writer, 500, "Error creating token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		ResponseWithError(writer, 500, "Error creating refresh token")
		return
	}

	err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		ResponseWithError(writer, 500, "Error writing token to database")
		return
	}

	ResponseWithJson(writer, 200, LoginUser{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email.String,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
