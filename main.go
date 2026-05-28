/*
	cmd to run and build server: go build -o out && ./out
	link to localhost:  http://localhost:8080/app/
*/

package main

import (
	"chirpy/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits	atomic.Int32
	db	*database.Queries
}

func main(){
	const filepathRoot = "."
	const port = "8080"
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db: dbQueries,
	}
	
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