package main

import (
	"log"
	"net/http"
	"time"
	"art/server"
)

func main() {
	//serving static files from ./public html/css
	fs := http.FileServer(http.Dir("./public"))
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// register HTTP handlers for different endpoints:
	// "/decoder" handles encoding/decoding POST requests
	mux.HandleFunc("/decoder", server.CodecHandler)
	// "/cypher" handles cypher POST requests
	mux.HandleFunc("/cypher", server.CypherHandler)
	// "/" servers the main index page (GET requests)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// returns 404 for any path other than "/"
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		server.IndexHandler(w, r)
	})

	// Creating rate limiter allowing 5 requests per 5 seconds per user.
	r1 := server.NewRateLimiter(5, 5*time.Second)

	// Wraps the mux router with the rate limiter middleware.
	handler := r1.Middleware(mux)

	// starts HTTP server on port 8080 with loogging
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err) // logs fatal error and stop if server fails to start.
	}
	

}

