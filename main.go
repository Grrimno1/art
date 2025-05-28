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
	mux.Handle("/static", http.StripPrefix("/static/", fs))

	//endpoints for /decoder and /cypher (POST) and "/" (GET)
	mux.HandleFunc("/decoder", server.CodecHandler)
	mux.HandleFunc("/cypher", server.CypherHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		server.IndexHandler(w, r)
	})

	r1 := server.NewRateLimiter(5, 5*time.Second)

	handler := r1.Middleware(mux)


	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}
	

}

