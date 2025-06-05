package server

import (
	"art/functions"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)
// HistoryEntry keeps track of each user's encode/decode operation for display in history.
type HistoryEntry struct {
	Timestamp 	string
	Action		string
	Input		string
	Result		string
}

var (
	history 	[]HistoryEntry 	// stores all recent history entries.
	historyMutex sync.Mutex 	//protects concurrent access to the history slice.
)

// constants used for action types and error messages.
const (
	actionEncode 		= "encode"
	actionDecode 		= "decode"
	errorString			= "Error\n"
)


/*
	handling /decoder POST requests for encoding and decoding.
	- Accepts only POST requests
	- Parses form inputs: encodeInput, decodeInput, and action.
	- Validation through validateInputs().
	- based on action -> calls processEncoding() or processDecoding().
	- success -> updates CombinedDataPage with results and status.
	- operation is saved in history.
	- Locks and copies history for template rendering.
	- Renders result page with the updated data.
*/

func CodecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Method Not Allowed: received %s, only POST allowed", r.Method)
		http.Error(w, MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	// prepare data struct for rendering template response
	data := CombinedPageData{
		Section: "decoder",
	}

	// parse POST form data
	if err := r.ParseForm(); err != nil {
			log.Printf("%s: %v", MsgFailedToParseForm, err)
			respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, MsgFailedToParseForm), &data)
			return
		}
		
		// normalize newlines for consistent processing
		rawDecodeInput := normalizeNewLines(r.FormValue("decodeInput"))
		rawEncodeInput := normalizeNewLines(r.FormValue("encodeInput"))
		action := r.FormValue("action")

		// validates inputs and gets any errors
		errMsg, statusType, statusCode, decodeInput, encodeInput := validateInputs(action, rawDecodeInput, rawEncodeInput)

		// if validation failed, render template with error message
		if errMsg != "" {
			data.StatusMessage = errMsg
			data.StatusCode = statusCode
			data.StatusType = statusType
			data.DecodeInput = decodeInput
			data.EncodeInput = encodeInput

			w.WriteHeader(statusCode)
			renderTemplate(w, data)
			return
		}
		
		// store raw inputs in data struct for rendering
		data.DecodeInput = rawDecodeInput
		data.EncodeInput = rawEncodeInput

		// process encoding or decoding based on action.
		switch action {
		case actionEncode:
			result, err := processEncoding(data.EncodeInput)
			if err != nil {
				// clears DecodeInput and respond with error on failure
				data.DecodeInput = ""
				respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, MsgMalformedInput), &data)
                return
			}
			data.DecodeInput = result
			data.StatusCode = http.StatusAccepted
			data.StatusType = statusSuccess
			data.StatusMessage = formatStatusMessage(http.StatusAccepted, MsgSuccessfullyEncoded)
			w.WriteHeader(http.StatusAccepted)
				
			saveHistory(actionEncode, data.EncodeInput, result)
			

		case actionDecode:
			result, err := processDecoding(data.DecodeInput)
			if err != nil {
				// clears EncodeInput and respond with error on failure
				data.EncodeInput = ""
				respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, MsgMalformedInput), &data)
                return
			}
			data.EncodeInput = result
			data.StatusCode = http.StatusAccepted
			data.StatusType = statusSuccess
			data.StatusMessage = formatStatusMessage(http.StatusAccepted, MsgSuccessfullyDecoded)
			w.WriteHeader(http.StatusAccepted)

			saveHistory(actionDecode, data.DecodeInput, result)
			

		default:
			// invalid action sent by user
			respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, MsgInvalidAction), &data)
			return
		}
	
	// calculate the number of lines for dynamic text area sizing
	data.LineCount = max(countLines(data.DecodeInput), countLines(data.EncodeInput))
	
	// safely copy history slice under lock to prevent concurrent access issues.
	historyMutex.Lock()
	data.History = make([]HistoryEntry, len(history))
	copy(data.History, history)
	historyMutex.Unlock()

	// render the updated page with encoding/decoding results and history.
	renderTemplate(w, data)
}

// processEncoding calls the encoding function and handles errors
func processEncoding(input string) (string, error) {
	result := functions.EncodeString(input, false)
	if result == errorString {
		return "", errors.New(MsgMalformedInput)
	}
	return result, nil
}

// same as above, but for decoding.
func processDecoding(input string) (string, error) {
	result := functions.DecodeString(input, false)
	if result == errorString {
		return "", errors.New(MsgMalformedInput)
	}
	return result, nil
}

/*
	validateInputs ensures that user inputs are valid:
		- Checks configuration sanity (MaxReturnLength vs MaxInputLength)
		- Requires at least one input field to be filled
		- validates length limits for encoding and decoding inputs
		- returns error message, status code, and sanitized inputs for rendering
*/
func validateInputs(action, rawDecodeInput, rawEncodeInput string) (errMsg, statusType string, statusCode int, decodeInput, encodeInput string) {
	if MaxReturnLength > MaxInputLength {
		log.Println("Misconfigured: maxReturnLen exceed maxInputLength")
		return formatStatusMessage(http.StatusInternalServerError, MsgInternalServerError), StatusError, http.StatusInternalServerError, "", ""
	}
	if rawDecodeInput == "" && rawEncodeInput == "" {
		return formatStatusMessage(http.StatusOK, MsgPleaseEnterText), StatusInfo, http.StatusOK, "", ""
	}
	if action == actionEncode && inputExceedsLimit(rawEncodeInput, MaxInputLength) {
		return formatStatusMessage(http.StatusRequestEntityTooLarge, MsgInputTooLong),
		 StatusError, http.StatusRequestEntityTooLarge, "", rawEncodeInput[:MaxReturnLength] + "..."
	}
	if action == actionDecode && decodedExceedsLimit(rawDecodeInput, MaxInputLength) && 
		!inputExceedsLimit(rawDecodeInput, MaxInputLength) {
		return formatStatusMessage(http.StatusUnprocessableEntity, MsgResultTooLong), 
		StatusError, http.StatusUnprocessableEntity, rawDecodeInput, ""
		}
	if action == actionDecode && inputExceedsLimit(rawDecodeInput, MaxInputLength) {
		return formatStatusMessage(http.StatusRequestEntityTooLarge, MsgInputTooLong),
		StatusError, http.StatusRequestEntityTooLarge, rawDecodeInput[:MaxReturnLength] + "...", ""
	}
	return "", "", 0, rawDecodeInput, rawEncodeInput
}
// saveHistory safely appends a new encode/decode operation to the history slice
func saveHistory(action, input, result string) {
	entry := HistoryEntry {
		Timestamp: 	time.Now().Format("January 2, 15:04"),
		Action: 	action,
		Input:		input,
		Result:		result,
	}
	// lock before modifying shared history slice to prevent race conditions
	historyMutex.Lock()
	defer historyMutex.Unlock()

	// insert new entry at the beginning of the slice.
	history = append([]HistoryEntry{entry}, history...)

	// track history to last MaxHistoryEntries entries to limit size
	if len(history) > MaxHistoryEntries {
		history = history[:MaxHistoryEntries]
	}
}

