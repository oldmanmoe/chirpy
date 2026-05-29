package main

import (
	"encoding/json"
	"net/http"
)

	type chirpRequest struct {
		Body   string	`json:"body"`	
	}

	type successResponse struct {
		Valid	bool	`json:"valid"`
	}

	type cleanResponse struct {
		Cleaned_Body	string	`json:"cleaned_body"`	
	}

func chirpCharLimitHandler(w http.ResponseWriter, req *http.Request) {
	var chirpResp chirpRequest
	
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&chirpResp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	
	const maxChirpLength = 140
	if len(chirpResp.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long" )
		return 
	}

	cleanChirpStr := cleanBadWords(chirpResp.Body)
	cleanChirp := cleanResponse{
		Cleaned_Body: cleanChirpStr,
	}

	respondWithJSON(w, http.StatusOK, cleanChirp)
	

	

}

