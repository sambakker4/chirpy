package main

import (
	"encoding/json"
	"net/http"
	"github.com/sambakker4/chirpy/internal/auth"

	"github.com/google/uuid"
)

type WebHook struct {
	Event string `json:"event"`
	Data struct{
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg apiConfig) HandleWebHook(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()	

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		ResponseWithError(writer, 401, "Error retrieving apiKey from header")
		return
	}

	if apiKey != cfg.apiKey {
		ResponseWithError(writer, 401, "Error, apikey not valid")
		return
	}

	var hook WebHook
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&hook)

	if err != nil {
		ResponseWithError(writer, 400, "Error decoding json")
		return
	}

	switch hook.Event {
		case "user.upgraded":
			UpgradeUser(&cfg, writer, hook, req)
			return
	}

	ResponseWithError(writer, 204, "event not supported")
}

func UpgradeUser(cfg *apiConfig, writer http.ResponseWriter, webHook WebHook, req *http.Request){
	_, err := cfg.db.SetUserChirpyRed(req.Context(), webHook.Data.UserID)
	if err != nil {
		ResponseWithError(writer, 404, "user not found")
		return
	}

	ResponseWithJson(writer, 204, struct{}{})
}
