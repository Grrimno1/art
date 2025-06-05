package server

// helper functions for data validation and template rendering

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"regexp"
	"strconv"
	"strings"
)
/* Predefined messages and constants for error handling and status reporting.
   Constants to improve maintainability.
*/
const (
	MsgMalformedInput		= "malformed input"
	MsgInvalidAction		= "invalid action"
	MsgFailedToParseForm	= "failed to parse form"
	MsgPleaseEnterText		= "please enter text to encode or decode"
	MsgInputEmpty			= "input cannot be empty"
	MsgXOREmpty				= "XOR key cannot be empty"
	MsgInputTooLong			= "input is too long, maximum length is 10,000 characters"
	MsgResultTooLong		= "result exceeds maximum length of 10,000 characters"
	MsgSuccessfullyEncoded 	= "successfully encoded"
	MsgSuccessfullyDecoded	= "successfully decoded"
	MsgInternalServerError 	= "internal server error"
	MsgMethodNotAllowed 	= "method not allowed"

	StatusInfo 				= "info"
	StatusError				= "error"
	statusSuccess			= "success"

	// limits for input and output lengths to prevent performance issues or abuse
	MaxInputLength 		= 10000
	MaxReturnLength 	= 1000
	MaxHistoryEntries 	= 20
	maxKeyLength 		= 256 // max length for XOR key
)

// parse and store the main HTML template used for rendering pages.
var tmpl = template.Must(template.ParseFiles("public/index.html"))

/* 
	CombinedPageData holds all the dynamic data passed into the HTML template.
	it covers both "art" and "cypher" sections of the app.
*/
type CombinedPageData struct {
	Section			string // specifies whether we're in 'art' or 'cypher' section

	// fields for the art encoder/decoder page
	DecodeInput		string
	EncodeInput		string
	StatusCode		int
	StatusType		string
	StatusMessage	string
	LineCount		int
	History			[]HistoryEntry

	// fields for the cypher page
	Mode			string
	Input			string
	Result			string
	Key				string
	CypherHistory	[]CypherHistoryEntry
}
// normalizeNewLines converts windows-style CRLF line endings ("\r\n")
// to Unix-style LF("\n") for consistent text processing.
func normalizeNewLines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
// formatStatusMessage creates a user-friendly status message that
// includes the HTTP status code, its text description, and a custom message.
func formatStatusMessage(code int, msg string) string {
	return fmt.Sprintf("%d %s: %s", code, http.StatusText(code), msg)
}

// respondWithError sets error details in the page data, writes the HTTP status code,
// and renders the error page template for the user.
func respondWithError(w http.ResponseWriter, code int, msg string, data *CombinedPageData) {
	data.StatusCode = code
	data.StatusType = StatusError //calling const from artCodec.go
	data.StatusMessage = msg
	w.WriteHeader(code)
	renderTemplate(w, *data)
}

// countLines counts the number of lines in a string (based on '\n')
// ensures a minimum of 4 lines for better textarea sizing in the UI.
func countLines(s string) int {
	lines := strings.Count(s, "\n") + 1
	if lines < 4 {
		return 4
	}
	return lines
}

// max returns the greater of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// decodedExceedsLimits estimates the length of the decoded text for
// inputs that use the custom compression format like '[count pattern]'
// Returns true if the decoded length would exceed the specified limit.
// This helps avoid decoding very large inputs that could cause performance issues.
func decodedExceedsLimit(input string, limit int) bool {
	pattern := regexp.MustCompile(`\[(\d+)\s([^\[\]]+)\]`)
	matches := pattern.FindAllStringSubmatchIndex(input, -1)

	decodedLen := 0
	prevEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]		// overall match indices
		countStr := input[match[2]:match[3]]	//repetition count as string
		patternStr := input[match[4]:match[5]]	// pattern to repeat

		// add length of literal text before this encoded chunk
		decodedLen += len(input[prevEnd:start])

		count, err := strconv.Atoi(countStr)	// convert count string to int
		if err != nil {
			continue // skip if count is invalid
		}
		decodedLen += count * len(patternStr)	// add length of repeated pattern * count

		// too long, return true.
		if decodedLen > limit {					
			return true
		}

		prevEnd = end 	// move past this match
		

	}

	// add any remaining literal text after the last encoded chunk
	decodedLen += len(input[prevEnd:])

	return decodedLen > limit
}
// inputExceedsLimit checks if the raw input string exceeds the maximum allowed length.
func inputExceedsLimit(input string, limit int) bool {
	return len(input) > limit
}
// renderTemplate executes the main HTML template with the provided data,
// sending the fully rendered page to the user's browser.
// Logs an error and sends a 500 error if rendering fails.
func renderTemplate(w http.ResponseWriter, data CombinedPageData) {
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}