package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
}

func main(){
	/*
	cmd to run and build server: go build -o out && ./out
	link to localhost:  http://localhost:8080/app/
	*/	
	
	const filepathRoot = "."
	const port = "8080"
	
	apiCfg := apiConfig{}
	handlerFileServer := http.FileServer(http.Dir(filepathRoot))
	
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(handlerFileServer)))
	
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/validate_chirp", chirpCharLimitHandler)

	srv := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}