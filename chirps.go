package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	writer.Header().Set("Content-Type", "application/json")
		
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding json: %s", err)
		writer.WriteHeader(500)
		return

	}

	if len(params.Body) > 140 {
		writer.WriteHeader(400)
		resp, _ := json.Marshal(returnError{Error: "Chirp is too long"})
		writer.WriteHeader(400)
		writer.Write(resp)
		return
	}
	
	resp, _ := json.Marshal(returnVals{Valid: true})
	writer.WriteHeader(200)
	writer.Write(resp)
}
