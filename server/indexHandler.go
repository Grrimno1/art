package server

import (
	"net/http"
	"log"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//validate for method that it is GET
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data := CombinedPageData{
		Section:		"art",
		DecodeInput:	"",
		EncodeInput: 	"",
		StatusCode:		http.StatusOK,
		StatusType:		"success",
		StatusMessage:  "200 Ok: welcome",
		LineCount:		4,	
	}
	
	
	if err := tmpl.Execute(w, data); err != nil {
		// for stdout
		log.Printf("IndexHandler: template execution error: %v", err)
		// log error server-side
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}