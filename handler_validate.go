package main

import (
	"encoding/json"
	"net/http"
)

func chirpCharLimitHandler(w http.ResponseWriter, req *http.Request) {

	type chirpResponse struct {
		Body   string	`json:"body"`	
	}

	type successResponse struct {
		Valid	bool	`json:"valid"`
	}


	decoder := json.NewDecoder(req.Body)
	chirpResp := chirpResponse{}
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

	resp := successResponse{
		Valid: true,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

