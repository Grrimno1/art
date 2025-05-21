package server

import (
	"art/functions"
	"net/http"
)

//Simple functions require simple file.


/*
	Handling /decoder POST requests.
	Body should contain data as x-www-form-urlencoded with key "input"
*/
func DecoderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	input := r.FormValue("input")
	if input == "" {
		http.Error(w, "Input cannot be empty", http.StatusBadRequest)
		return
	}

	output := functions.DecodeString(input, true)
	if output == "Error\n" {
		http.Error(w, "Malformed input", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(output))
}


/*
	Handling /encoder POST requests.
	Body should contain data as x-www-form-urlencoded with key "input"
*/
func EncoderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	input := r.FormValue("input")
	if input == "" {
		http.Error(w, "Input cannot be empty", http.StatusBadRequest)
		return
	}

	output := functions.EncodeString(input, true)
	if output == "Error\n" {
		http.Error(w, "Invalid input for encoding", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(output))
}