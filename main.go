package main

import (
	"log"
	"net/http"
)

func main(){
	//cmd to run and build server: go build -o out && ./out
	
	const filepathRoot = "."
	const port = "8080"
	
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	srv := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}
	

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	

}