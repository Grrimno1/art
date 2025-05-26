package server

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := CombinedPageData{
		Section:		"art",
		DecodeInput:	"",
		EncodeInput: 	"",
		StatusCode:		http.StatusOK,
		StatusType:		"success",
		StatusMessage:  "200 Ok: welcome",
		LineCount:		4,	
	}
	
	err := tmpl.Execute(w, data)
	if err != nil {
		// log error server-side
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		// return is good here, keep it
		return
	}
}