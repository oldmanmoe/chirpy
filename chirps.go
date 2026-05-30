package main

import (
	"chirpy/internal/database"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

	type chirpRequest struct {
		Body   	string		`json:"body"`
		UserID	uuid.UUID	`json:"user_id"`

	}

	type successResponse struct {
		Valid	bool	`json:"valid"`
	}

	type Chirp struct {
		ID			uuid.UUID	`json:"id"`
		CreatedAt	time.Time	`json:"created_at"`
		UpdatedAt	time.Time	`json:"updated_at"`
		Body		string		`json:"body"`
		UserID		uuid.UUID	`json:"user_id"`
	}

func(cfg *apiConfig) chirpCharLimitHandler(w http.ResponseWriter, req *http.Request) {
	var chirpReq chirpRequest
	
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	
	const maxChirpLength = 140
	if len(chirpReq.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long" )
		return 
	}

	cleanChirpStr := cleanBadWords(chirpReq.Body)
	

	
	chirpInfo, err := cfg.db.CreateChirp(context.Background(),
	database.CreateChirpParams{
		Body: cleanChirpStr,
		UserID: chirpReq.UserID,
		})

	if err != nil {
		log.Fatal(err)
		respondWithError(w, http.StatusNotAcceptable, "something went wrong creating chirp")
		return
	}

	finalChirp := Chirp{
		ID: chirpInfo.ID,
		CreatedAt: chirpInfo.CreatedAt,
		UpdatedAt: chirpInfo.UpdatedAt,
		Body: chirpInfo.Body,
		UserID: chirpReq.UserID,
	}


	respondWithJSON(w, http.StatusCreated, finalChirp)
	

	

}

