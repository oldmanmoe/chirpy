package main

import (
	internal "chirpy/internal/auth"
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

func(cfg *apiConfig) chirpRequestHandler(w http.ResponseWriter, r *http.Request) {
	var chirpReq chirpRequest
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	authToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required")
		return
	}

	reqUser, err := internal.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized, authentication required!")
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
		UserID: reqUser,
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
		UserID: reqUser,
	}

	respondWithJSON(w, http.StatusCreated, finalChirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	var result []Chirp

	allChirps, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to fetch all chirps")
		return
	}

	for _, chirp := range allChirps {
		result = append(result,
			Chirp{
				ID: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserID: chirp.UserID,

			})
		}

	respondWithJSON(w, http.StatusOK, result)
}

func (cfg *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {

	rawChirpID := r.PathValue("chirpId")
	chirpID, err := uuid.Parse(rawChirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong parsing rawChirpID")
		return
	}

	chirpInfo, err := cfg.db.GetChirp(context.Background(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	result := Chirp{
		ID: chirpInfo.ID,
		CreatedAt: chirpInfo.CreatedAt,
		UpdatedAt: chirpInfo.UpdatedAt,
		Body: chirpInfo.Body,
		UserID: chirpInfo.UserID,
	}

	respondWithJSON(w, http.StatusOK, result)
}

