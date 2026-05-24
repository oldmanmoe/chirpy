package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type errorResponse struct {
		Error	string	`json:"error"`
	}

	respBody := errorResponse {
		Error: msg,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	payloadResp, err:= json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshallin JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payloadResp)

}