package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}


func main() {
	const port = "8080"
	const filePathRoot = "./app/"

	mux := http.NewServeMux()
	config := &apiConfig{}
	handler := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/", http.StripPrefix("/app", config.middlewareMetricsInc(handler)))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", config.getFileserverHits)
	mux.HandleFunc("POST /reset", config.resetFileserverHits)

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
