package server

import (
	"art/functions"
	"html/template"
	"net/http"
)

type PageData struct {
	DecodeInput string
	EncodeInput string
}

var tmpl = template.Must(template.ParseFiles("public/index.html"))

/* 
	handling /decoder POST requests for both encoding and decoding
*/

func CodecHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}
	
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return
		}

		data.DecodeInput = r.FormValue("decodeInput")
		data.EncodeInput = r.FormValue("encodeInput")

		action := r.FormValue("action")

		switch action {
		case "encode":
			//Encode input to decodeOutput (swap)
			result := functions.EncodeString(data.EncodeInput, true)
			if result == "Error\n" {
				http.Error(w, "Malformed input", http.StatusBadRequest)
				return
			}
			data.DecodeInput = result // shows encoded result in decode textarea.
		
		case "decode":
			//Decode input to encodeOutput (swap)
			result := functions.DecodeString(data.DecodeInput, true)
			if result == "Error\n" {
				http.Error(w, "Malformed input", http.StatusBadRequest)
				return
			}
			data.EncodeInput = result //shows decoded result in encode textarea.

		default: 
			http.Error(w, "Invalid action: use 'encode' or 'decode'", http.StatusBadRequest)
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
}