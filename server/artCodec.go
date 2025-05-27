package server

import (
	"art/functions"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)
type HistoryEntry struct {
	Timestamp 	string
	Action		string
	Input		string
	Result		string
}

var (
	history 	[]HistoryEntry
	historyMutex sync.Mutex
	tmpl = template.Must(template.ParseFiles("public/index.html"))
)

const (
	maxInputLength = 10000
	maxReturnLen = 1000
)


/*
	handling /decoder POST requests for both encoding and decoding
*/

func CodecHandler(w http.ResponseWriter, r *http.Request) {
	//allowing only POST requests from clientside.
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		log.Printf("Method Not Allowed: received %s, only POST allowed", r.Method)
		return
	}
	data := CombinedPageData{
		Section: "decoder",
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, "Failed to parse form"), &data)
			log.Printf("Failed to parse form")
			return
		}
		
		//parsing inputs from form.
		rawDecodeInput := strings.ReplaceAll(r.FormValue("decodeInput"), "\r\n", "\n")
		rawEncodeInput := strings.ReplaceAll(r.FormValue("encodeInput"), "\r\n", "\n")
		action := r.FormValue("action")

		//help variables for validation
		
		//validating data from form
		errMsg, statusType, statusCode, decodeInput, encodeInput := validateInputs(action, rawDecodeInput, rawEncodeInput)

		if errMsg != "" {
			data.StatusMessage = errMsg
			data.StatusCode = statusCode
			data.StatusType = statusType
			data.DecodeInput = decodeInput
			data.EncodeInput = encodeInput

			w.WriteHeader(statusCode)
			_ = tmpl.Execute(w, data)
			return
		}
		
		data.DecodeInput = rawDecodeInput
		data.EncodeInput = rawEncodeInput

		switch action {
		case "encode":
			result := functions.EncodeString(data.EncodeInput, false)
			if result == "Error\n" {
				respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, "Malformed input"), &data)
                return
			} else {
				data.DecodeInput = result
				data.StatusCode = http.StatusAccepted
				data.StatusType = "success"
				data.StatusMessage = formatStatusMessage(http.StatusAccepted, "Successfully encoded.")
				w.WriteHeader(http.StatusAccepted)
				
				//saving input in history.
				saveHistory("encode", data.EncodeInput, result)
			}

		case "decode":
			result := functions.DecodeString(data.DecodeInput, false)
			if result == "Error\n" {
				respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, "Malformed input"), &data)
                return
			} else {
				data.EncodeInput = result
				data.StatusCode = http.StatusAccepted
				data.StatusType = "success"
				data.StatusMessage = formatStatusMessage(http.StatusAccepted, "Successfully decoded.")
				w.WriteHeader(http.StatusAccepted)

				//saving input in history.
				saveHistory("decode", data.DecodeInput, result)
			}

		default:
			respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, "Invalid action."), &data)
			return
		}
	}
	
	data.LineCount = max(CountLines(data.DecodeInput), CountLines(data.EncodeInput))
	
	//safely reading history before assigning template data.
	historyMutex.Lock()
	data.History = make([]HistoryEntry, len(history))
	copy(data.History, history)
	historyMutex.Unlock()

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		log.Printf("template execution error: %v", err)
	}
}

func validateInputs(action, rawDecodeInput, rawEncodeInput string) (errMsg, statusType string, statusCode int, decodeInput, encodeInput string) {
	if maxReturnLen > maxInputLength {
		log.Println("Misconfigured: maxReturnLen exceed maxInputLength")
		return formatStatusMessage(http.StatusInternalServerError, "Internal server error."), "error", http.StatusInternalServerError, "", ""
	}
	if rawDecodeInput == "" && rawEncodeInput == "" {
		return formatStatusMessage(http.StatusOK, "please enter text to encode or decode."), "info", http.StatusOK, "", ""
	}
	if action == "encode" && inputExceedsLimit(rawEncodeInput, maxInputLength) {
		return formatStatusMessage(http.StatusRequestEntityTooLarge, "input is too long, Maximum length is 10,000 characters"),
		 "error", http.StatusRequestEntityTooLarge, "", rawEncodeInput[:maxReturnLen] + "..."
	}
	if action == "decode" && decodedExceedsLimit(rawDecodeInput, maxInputLength) && 
		!inputExceedsLimit(rawDecodeInput, maxInputLength) {
		return formatStatusMessage(http.StatusUnprocessableEntity, "Result exceeds Maximum length of 10,000 characters"), 
		"error", http.StatusUnprocessableEntity, rawDecodeInput, ""
		}
	if action == "decode" && inputExceedsLimit(rawDecodeInput, maxInputLength) {
		return formatStatusMessage(http.StatusRequestEntityTooLarge, "input is too long, Maximum length is 10,000 characters"),
		"error", http.StatusRequestEntityTooLarge, rawDecodeInput[:maxReturnLen] + "...", ""
	}
	return "", "", 0, rawDecodeInput, rawEncodeInput
}

func formatStatusMessage(code int, msg string) string {
	return fmt.Sprintf("%d %s: %s", code, http.StatusText(code), msg)
}
func respondWithError(w http.ResponseWriter, code int, msg string, data *CombinedPageData) {
	data.StatusCode = code
	data.StatusType = "error"
	data.StatusMessage = msg
	w.WriteHeader(code)
	_ = tmpl.Execute(w, data)
}

func saveHistory(action, input, result string) {
	entry := HistoryEntry {
		Timestamp: 	time.Now().Format("02.01.06 15:04"),
		Action: 	action,
		Input:		input,
		Result:		result,
	}
	//locking so we can append user data without issues. This is to avoid duplicate entries.
	historyMutex.Lock()
	defer historyMutex.Unlock()

	history = append(history, entry)
	//keep track of the last 20 entries.
	if len(history) > 20 {
		history = history[len(history)-20:]
	}
}

//counting the lines of input so we can adjust textareas in frontend for better user experience.
func CountLines(s string) int {
	lines := strings.Count(s, "\n") + 1
	if lines < 4 {
		return 4
	}
	return lines
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
