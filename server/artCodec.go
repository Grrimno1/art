package server

import (
	"art/functions"
	"net/http"
)

/* 
	handling /decoder POST requests for both encoding and decoding
*/

func CodecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}
	
	mode := r.FormValue("mode")
	input := r.FormValue("input")

	var output string
	switch mode {
	case "encode":
		output = functions.EncodeString(input, true)
	case "decode":
		output = functions.DecodeString(input, true)
	default:
		http.Error(w, "Invalid mode: use 'encode' or 'decode'", http.StatusBadRequest)
	}

	if output == "Error\n" {
		http.Error(w, "Malformed input", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(output))
}