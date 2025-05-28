package server

import (
	"art/functions"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)
// struct to keep history of user encode/decode inputs.
type HistoryEntry struct {
	Timestamp 	string
	Action		string
	Input		string
	Result		string
}

var (
	history 	[]HistoryEntry
	historyMutex sync.Mutex
)

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
	data := CombinedPageData{
		Section: "decoder",
	}

	if err := r.ParseForm(); err != nil {
			log.Printf("%s: %v", MsgFailedToParseForm, err)
			respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, MsgFailedToParseForm), &data)
			return
		}
		
		rawDecodeInput := normalizeNewLines(r.FormValue("decodeInput"))
		rawEncodeInput := normalizeNewLines(r.FormValue("encodeInput"))
		action := r.FormValue("action")

		//validating data from form
		errMsg, statusType, statusCode, decodeInput, encodeInput := validateInputs(action, rawDecodeInput, rawEncodeInput)

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
		
		data.DecodeInput = rawDecodeInput
		data.EncodeInput = rawEncodeInput

		switch action {
		case actionEncode:
			result, err := processEncoding(data.EncodeInput)
			if err != nil {
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
			respondWithError(w, http.StatusBadRequest, formatStatusMessage(http.StatusBadRequest, MsgInvalidAction), &data)
			return
		}
	
	// Dynamic calculation for textareas; increase @ frontend if DecodeInput or EncodeInput row > 4
	data.LineCount = max(countLines(data.DecodeInput), countLines(data.EncodeInput))
	
	//safely copying history slice under mutex lock
	historyMutex.Lock()
	data.History = make([]HistoryEntry, len(history))
	copy(data.History, history)
	historyMutex.Unlock()

	renderTemplate(w, data)
}

// Encodes data using part1 functions
func processEncoding(input string) (string, error) {
	result := functions.EncodeString(input, false)
	if result == errorString {
		return "", errors.New(MsgMalformedInput)
	}
	return result, nil
}

// decodes data using part1 functions
func processDecoding(input string) (string, error) {
	result := functions.DecodeString(input, false)
	if result == errorString {
		return "", errors.New(MsgMalformedInput)
	}
	return result, nil
}

/*input validation; checks for input length limits. 
  added this to improve security of this webserver.
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
// appends new encode/decode record to history with a timestamp.
func saveHistory(action, input, result string) {
	entry := HistoryEntry {
		Timestamp: 	time.Now().Format("2006-01-02 15:04"),
		Action: 	action,
		Input:		input,
		Result:		result,
	}
	//locking so we can append user data without issues. This is to avoid duplicate entries.
	historyMutex.Lock()
	defer historyMutex.Unlock()

	history = append(history, entry)
	//keep track of the last 20 entries.
	if len(history) > MaxHistoryEntries {
		history = history[len(history) - MaxHistoryEntries:]
	}
}

