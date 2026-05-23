package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
}

func readinessHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main(){
	//cmd to run and build server: go build -o out && ./out
	
	const filepathRoot = "."
	const port = "8080"
	
	cfg := apiConfig{}
	handlerFileServer := http.FileServer(http.Dir(filepathRoot))
	
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(handlerFileServer)))
	
	mux.HandleFunc("/healthz", readinessHandler)
	mux.HandleFunc("/metrics", cfg.serverHitsHandler)
	mux.HandleFunc("/reset", cfg.resetHitsHandler)

	srv := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}