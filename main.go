package main

import (
	"log"
	"net/http"
	"art/server"
)

func main() {
	//for server static files.
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs) //index

	// endpoint for decoder & encoder
	http.HandleFunc("/decoder", server.CodecHandler)
	// endpoint for cypher (XOR/rot13)
	http.HandleFunc("/cypher", server.CypherHandler)
	// initializes server at port 8080
	log.Println("Server starting on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}


