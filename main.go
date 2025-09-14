package main

import (
	// "fmt"

	"database/sql"
	"log"
	"net/http"

	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"Chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
	}

	// 1. Create a new http.ServeMux and register a handler
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("GET /admin/metrics", apiCfg.statsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetStatsHandler)

	// 2. Create a new http.Server struct and assign the mux as its handler
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Use the http.FileServer to serve files from the current directory
	//fileServer := http.FileServer(http.Dir("."))
	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	// Use the mux's .Handle() method to add a handler for the root path
	// Strip the /app prefix from the request path before passing it to the fileserver handler
	mux.Handle("/app/", fileServer)

	// Also handle the no-trailing-slash case so /app serves index (or redirects):
	// Redirect /app ==> /app/
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/app/", http.StatusMovedPermanently)
	})

	// 3. Use the server's ListenAndServe method
	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
