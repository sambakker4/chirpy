package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) getFileserverHits(writer http.ResponseWriter, req *http.Request) {
	hits := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	writer.Write([]byte(hits))
}

func (cfg *apiConfig) resetFileserverHits(writer http.ResponseWriter, req *http.Request) {
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
