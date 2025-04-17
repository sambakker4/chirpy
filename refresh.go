package main

import (
	"net/http"
	"time"

	"github.com/sambakker4/chirpy/internal/auth"
)

type TokenResponse struct {
	Token string `json:"token"`
}


func (cfg apiConfig) Refresh(writer http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		ResponseWithError(writer, 400, "Error retrieving refresh token from header")
		return
	}

	token, err := cfg.db.GetRefreshToken(req.Context(), tokenString)
	if err != nil {
		ResponseWithError(writer, 401, "Error retrieving token from database")
		return
	}

	if revoked, _ := token.RevokedAt.Value(); revoked != nil {
		ResponseWithError(writer, 401, "Token revoked")
	}

	if token.ExpiresAt.Before(time.Now()) {
		ResponseWithError(writer, 401, "Token expired")	
		return
	}

	user, err := cfg.db.GetUserByRefreshToken(req.Context(), token.Token)
	if err != nil {
		ResponseWithError(writer, 500, "Error retrieving user from database")
		return
	}
	
	responseToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, time.Hour)		
	if err != nil {
		ResponseWithError(writer, 500, "Error creating access token")
		return
	}

	ResponseWithJson(writer, 200, TokenResponse{
		Token: responseToken,
	})
}
