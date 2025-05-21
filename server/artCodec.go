package server

import (
	"art/functions"
	"io"
	"net/http"
)

func DecoderHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	input := string(body)

	output := functions.DecodeString(input, true)
	if output == "Error\n" {
		http.Error(w, "Malformed input", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted) //202
	w.Write([]byte(output))

}

func EncoderHandler (w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	input := string(body)

	output := functions.EncodeString(input, true)
	if output == "Error\n" {
		http.Error(w, "Invalid input for encoding", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(output))
}