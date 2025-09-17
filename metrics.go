package main

import (
	"fmt"
	"log"
	"net/http"
)

// to display the number of hits
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits.Load())))
}

// Middleware function (High Order Function) to count access to servers API
// WRAPPER or DECORATOR pattern
// Middleware is a way to wrap a handler with additional functionality.
// It is a common pattern in web applications that allows us to write DRY code.
// For example, we can write a middleware that logs every request to the server.
// We can then wrap our handler with this middleware and every request will be
// logged without us having to write the logging code in every handler.
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r) // call original
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
