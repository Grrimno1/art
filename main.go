package main

import (
	"log"
	"net/http"
	"art/server"
)

func main() {
	//server static files...
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs) //index

	// endpoint for decoder
	http.HandleFunc("/decoder", server.DecoderHandler)
	// endpoint for encoder
	http.HandleFunc("/encoder", server.EncoderHandler)

	log.Println("Server starting on: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}


