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
	db				*database.Queries
	platform 		string
	secret			string
	polkaKey		string
	
}

func main(){
	const filepathRoot = "."
	const port = "8080"
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		db: 		dbQueries,
		platform: platform,
		secret:	jwtSecret,
		polkaKey: polkaKey,
	}
	
	handlerFileServer := http.FileServer(http.Dir(filepathRoot))
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(handlerFileServer)))
	
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.handlerGetSingleChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.chirpRequestHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdatePassword)
	mux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.handlerDeleteSingleChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)
	

	srv := http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}