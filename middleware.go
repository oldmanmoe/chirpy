package main

import (
	"fmt"
	"net/http"
)


func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) serverHitsHandler(w http.ResponseWriter, req *http.Request) {
	result := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(result))
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) resetHitsHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Swap(0)
	w.WriteHeader(http.StatusOK)
}