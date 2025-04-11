package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)



func main() {
	const port = "8080"
	const filePathRoot = "./app/"

	mux := http.NewServeMux()
	config := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	handler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", config.getFileserverHits)
	mux.HandleFunc("POST /admin/reset", config.resetFileserverHits)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	fmt.Printf("Serving files from: %v on port: %v\n", filePathRoot, port)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
