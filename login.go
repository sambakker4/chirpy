package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/sambakker4/chirpy/internal/auth"
)

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg apiConfig) Login(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

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

	if err != nil {
		ResponseWithError(writer, 400, "Error hashing password")
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, loginInfo.Password)
	if err != nil {
		ResponseWithError(writer, 401, "incorrect email or password") 
		return
	}

	ResponseWithJson(writer, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email.String,
	})
}
