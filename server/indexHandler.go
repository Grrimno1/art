package server

import (
	"net/http"
)
// IndexHandler handles the request for the main (index) page.
// only accepts GET requests and renders the default "art" section with empty input fields.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// only allow GET requests to this handler. Reject anything else.
	if r.Method != http.MethodGet {
		http.Error(w, MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	/* 	Prepare the default data to be passed to the HTML template.
		this sets the initial state of the web form: 
	 		- section "art" will be visible.
			- no pre-filled input values.
			- a status message welcoming the user.
			- linecount is used to set the initial size of the input textarea.
	 */
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