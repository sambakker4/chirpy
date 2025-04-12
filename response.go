package main

import (
	"encoding/json"
	"net/http"
)

func ResponseWithError(writer http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}

	resp, _ := json.Marshal(returnError{Error: msg})
	writer.WriteHeader(code)
	writer.Write(resp)
}

func ResponseWithJson(writer http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		ResponseWithError(writer, 500, "Error marshaling return data")
		return
	}
	writer.WriteHeader(code)
	writer.Write(resp)
}
