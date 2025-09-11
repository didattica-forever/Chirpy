package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the homepage!")
}

func main() {
	// 1. Create a new http.ServeMux and register a handler
	mux := http.NewServeMux()
	//mux.HandleFunc("/", homeHandler)

	// 2. Create a new http.Server struct and assign the mux as its handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 3. Use the server's ListenAndServe method
	log.Println("Server starting on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}