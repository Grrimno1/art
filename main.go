package main

import (
	"log"
	"net/http"
	"art/server"
)

func main() {
	// serving static files from ./public html/css
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// end point for "/" GET's index.html
	http.HandleFunc("/", server.IndexHandler)

	// end point for /decoder and /cypher handles POST
	http.HandleFunc("/decoder", server.CodecHandler)
	http.HandleFunc("/cypher", server.CypherHandler)
	
	// starting server @ port 8080
	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}


}

