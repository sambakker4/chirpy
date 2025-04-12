package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"slices"
)

func validateChirp(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string`json:"cleaned_body"`
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

	message := params.Body

	if len(message) > 140 {
		ResponseWithError(writer, 400, "Chirp is too long")
		return
	}
	message = removeBadWords(message)
	ResponseWithJson(writer, 200, returnVals{CleanedBody: message})
}

func removeBadWords(str string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	words := strings.Split(str, " ")

	for i, word := range words {
		if idx := slices.Index(badWords, strings.ToLower(word)); idx != -1 {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
