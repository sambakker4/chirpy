package main

import (
	"net/http"
	"github.com/sambakker4/chirpy/internal/auth"
)



func (cfg apiConfig) RevokeToken(writer http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)	
	if err != nil {
		ResponseWithError(writer, 400, "Error reading token from header")
		return
	}

	err = cfg.db.RevokeToken(req.Context(), token)
	if err != nil {
		ResponseWithError(writer, 500, "Error revoking token")
		return
	}

	writer.WriteHeader(204)
	writer.Write([]byte(""))
}
