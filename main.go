package main

import (
	"log"
	"net/http"
	"art/server"
)

func main() {
	// serve static files from ./public/static (styles, etc)
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Use your server package to render template for "/"
	http.HandleFunc("/", server.IndexHandler)

	http.HandleFunc("/decoder", server.CodecHandler)
	http.HandleFunc("/cypher", server.CypherHandler)
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

