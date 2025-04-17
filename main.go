package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"database/sql"
	"github.com/sambakker4/chirpy/internal/database"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	tokenSecret := os.Getenv("TOKEN_SECRET")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	const port = "8080"
	const filePathRoot = "./app/"

	mux := http.NewServeMux()
	config := apiConfig{
		db: dbQueries,
		fileserverHits: atomic.Int32{},
		platform: platform,
		tokenSecret: tokenSecret,
	}

	handler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(handler)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", config.getFileserverHits)
	mux.HandleFunc("POST /admin/reset", config.resetFileserver)
	mux.HandleFunc("POST /api/chirps", config.CreateChirp)
	mux.HandleFunc("POST /api/users", config.CreateUser)
	mux.HandleFunc("GET /api/chirps", config.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.GetChirp)
	mux.HandleFunc("POST /api/login", config.Login)
	mux.HandleFunc("POST /api/refresh", config.Refresh)
	mux.HandleFunc("POST /api/revoke", config.RevokeToken)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Serving files from: %v on port: %v\n", filePathRoot, port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
