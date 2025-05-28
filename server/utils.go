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
	StatusInfo 				= "info"
	StatusError				= "error"

	// limits for validation
	MaxInputLength 		= 10000
	MaxReturnLength 	= 1000
	MaxHistoryEntries 	= 20
	maxKeyLength 		= 256 //for XOR
)
// for artCodec.go and cypherCodec.go
var tmpl = template.Must(template.ParseFiles("public/index.html"))

// Holds all data passed to HTML template.
// covers art & cypher.
type CombinedPageData struct {
	Section			string // art or cypher

	// art decoder/encoder fields
	DecodeInput		string
	EncodeInput		string
	StatusCode		int
	StatusType		string
	StatusMessage	string
	LineCount		int
	History			[]HistoryEntry

	//cypher fields
	Mode			string
	Input			string
	Result			string
	Key				string
	CypherHistory	[]CypherHistoryEntry
}
// Converts windows CRLF to Unix
func normalizeNewLines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}
// for user friendly status message.
func formatStatusMessage(code int, msg string) string {
	return fmt.Sprintf("%d %s: %s", code, http.StatusText(code), msg)
}

// sets error details in the CombinedPageData,
// writes the HTTP status code to the response and renders the template with the updated data.
func respondWithError(w http.ResponseWriter, code int, msg string, data *CombinedPageData) {
	data.StatusCode = code
	data.StatusType = StatusError //calling const from artCodec.go
	data.StatusMessage = msg
	w.WriteHeader(code)
	renderTemplate(w, *data)
}

// dynamic textareas for GUI
func countLines(s string) int {
	lines := strings.Count(s, "\n") + 1
	if lines < 4 {
		return 4
	}
	return lines
}

// increases textareas dynamically.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// validates the input to check if the result would exceed the set MaxInputLength if decoded.
func decodedExceedsLimit(input string, limit int) bool {
	pattern := regexp.MustCompile(`\[(\d+)\s([^\[\]]+)\]`)
	matches := pattern.FindAllStringSubmatchIndex(input, -1)

	decodedLen := 0
	prevEnd := 0

	for _, match := range matches {
		start, end := match[0], match[1]		// overall match indices
		countStr := input[match[2]:match[3]]	//repetition count as string
		patternStr := input[match[4]:match[5]]	// pattern to repeat

		decodedLen += len(input[prevEnd:start]) // add length of literal text before this match

		count, err := strconv.Atoi(countStr)	// convert count string to int
		if err != nil {
			continue // skip if count is invalid
		}
		decodedLen += count * len(patternStr)	// add length of repeated pattern * count

		// too long, return true.
		if decodedLen > limit {					// exit if limit exceeded
			return true
		}

		prevEnd = end 							// update prevEnd to after this match.
		

	}

	//adding remaining literal characters after last match
	decodedLen += len(input[prevEnd:])

	return decodedLen > limit
}
// checks if given raw input exceeds MaxInputLength
func inputExceedsLimit(input string, limit int) bool {
	return len(input) > limit
}
// executes HTML template with the provided data.
func renderTemplate(w http.ResponseWriter, data CombinedPageData) {
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}