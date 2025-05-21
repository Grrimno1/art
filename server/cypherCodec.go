package server

import (
	"io"
	"net/http"
	"art/functions"
	"strings"
)

func CypherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	mode := r.URL.Query().Get("mode")
	key := r.URL.Query().Get("key")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	data := string(body)

	var result string

	switch mode {
		case "xor":
			res, _ := functions.Xorify(data, key)
			if strings.HasPrefix(res, "Error") {
				http.Error(w, res, http.StatusBadRequest)
				return
			}
			result = res
	
		case "rot13":
			result = functions.Rot13ify(data)

		default:
			http.Error(w, "Invalid mode. Use 'xor' or 'rot13'", http.StatusBadRequest)
			return
		}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}