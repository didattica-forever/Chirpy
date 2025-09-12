package main

import (
	//"fmt"
	"log"
	"net/http"
)

// Readiness endpoints are commonly used by external systems
// to check if our server is ready to receive traffic.
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	// 1. Create a new http.ServeMux and register a handler
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", readinessHandler)

	// 2. Create a new http.Server struct and assign the mux as its handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Use the http.FileServer to serve files from the current directory
	//fileServer := http.FileServer(http.Dir("."))
	fileServer := http.FileServer(http.Dir("."))

	// Use the mux's .Handle() method to add a handler for the root path
	// Strip the /app prefix from the request path before passing it to the fileserver handler
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

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
