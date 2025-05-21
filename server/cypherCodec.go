package server

import (
	"net/http"
	"art/functions"
	"strings"
)
/* handles post requests for XOR and ROT13
	must be x-www-form-urlencoded and contain
		'mode' - determines whether to XOR or ROT13 the input
		'key' - for XOR
		'input' - data to encrypt/decrypt
*/
func CypherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	mode := r.FormValue("mode")
	key := r.FormValue("key")
	data := r.FormValue("input")

	if data == "" {
		http.Error(w, "Input cannot be empty", http.StatusBadRequest)
		return
	}

	var result string
	switch mode {
	case "xor":
		res, err := functions.Xorify(data, key)
		if err != nil {
			http.Error(w, "XOR error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if strings.HasPrefix(res, "Error") {
			http.Error(w, res, http.StatusBadRequest)
			return
		}
		result = res

	case "rot13":
		result = functions.Rot13ify(data)
	
	default:
		http.Error(w, "Invalid mode: Use 'xor' or 'rot13'", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}