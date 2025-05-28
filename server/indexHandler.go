package server

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//validate for method that it is GET
	if r.Method != http.MethodGet {
		http.Error(w, MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	data := CombinedPageData{
		Section:		"art",
		DecodeInput:	"",
		EncodeInput: 	"",
		StatusCode:		http.StatusOK,
		StatusType:		statusSuccess,
		StatusMessage:  formatStatusMessage(http.StatusOK, "welcome"),
		LineCount:		4,	
	}
	renderTemplate(w, data)
}