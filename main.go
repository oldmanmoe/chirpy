package main

import (
	"log"
	"net/http"
)

func main(){

	//code to run and build server: go build -o out && ./out
	
	mux := http.NewServeMux()

	srv := http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	

}