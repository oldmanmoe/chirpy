package main

import (
	"chirpy/internal/auth"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Event struct {
	Event	string `json:"event"`
	Data	struct {
		UserID	uuid.UUID	`json:"user_id"`
	} `json:"data"`
}


func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	var event Event

	reqAPI, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Something went wrong")
		return
	}

	if reqAPI != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required")
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&event)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if event.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeToChirpyRed(context.Background(), event.Data.UserID)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	} else if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}


	w.WriteHeader(http.StatusNoContent)
}
