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
	}

	handler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(handler)))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", config.getFileserverHits)
	mux.HandleFunc("POST /admin/reset", config.resetFileserver)
	mux.HandleFunc("POST /api/chirps", config.CreateChirp)
	mux.HandleFunc("POST /api/users", config.CreateUser)

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
