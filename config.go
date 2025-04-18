package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/sambakker4/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
	apiKey         string
}

func (cfg *apiConfig) getFileserverHits(writer http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`,
		cfg.fileserverHits.Load())
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte(hits))
}

func (cfg *apiConfig) resetFileserver(writer http.ResponseWriter, req *http.Request) {
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		log.Printf("Error reseting database: %v", err)
	}

	cfg.fileserverHits.Store(0)

	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte("Reset Successful"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(writer, req)
	})
}
