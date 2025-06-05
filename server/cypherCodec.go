package server

import (
	"net/http"
	"art/functions"
	"log"
	"sync"
	"time"
)
// CypherHistoryEntry stores details of each XOR or ROT13 opeation perfomed by the user.
type CypherHistoryEntry struct {
	Timestamp 	string
	Mode		string
	Key			string //optional: used for XOR
	Input		string
	Result		string
}

// operation identifiers.
const (
	xor		= "xor"
	rot13	= "rot13"
)
var (
	cypherHistory 		[]CypherHistoryEntry 	// stores recent cypher operations history
	cypherHistoryMutex	sync.Mutex 				// mutex to safely access cypherHistory slice concurrently
)

/* 
	handles post requests for XOR and ROT13
		must be x-www-form-urlencoded and contain
			- expects 'mode' (xor or rot13), 'key' (for xor) and 'input' (data to process) form values.
			- validates inputs for presence and length.
			- calls corresponding function for the requested mode.
			- records the operation in history.
			- return processed result or error status to the user.
*/
func CypherHandler(w http.ResponseWriter, r *http.Request) {
	// only POST requests are allowed; otherwise, return 405 error.
	if r.Method != http.MethodPost {
		log.Printf("cypherHandler: Method Not Allowed: received %s, only POST allowed", r.Method)
		http.Error(w, MsgMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}
	// prepare data struct for rendering template response
	data := CombinedPageData{
		Section: "cypher",
	}

	// extract form inputs
	mode := r.FormValue("mode")
	key := r.FormValue("key")
	rawInput := normalizeNewLines(r.FormValue("input"))

	// validate that input is not empty
	if rawInput == "" {
		log.Printf("cypherHandler: Empty input received")
		respondWithError(w, http.StatusBadRequest, MsgInputEmpty, &data)
		return
	}

	// if XOR mode selected, validate key is provided and not too long
	if mode == xor {
		if key == "" {
			log.Printf("cypherHandler: XOR mode selected but empty key provided")
			respondWithError(w, http.StatusBadRequest, MsgXOREmpty, &data)
			return
		}
		if inputExceedsLimit(key, maxKeyLength) {
			respondWithError(w, http.StatusRequestEntityTooLarge, formatStatusMessage(http.StatusRequestEntityTooLarge, "XOR key is too long"), &data)
			log.Printf("cypherHandler: XOR key too long: %d characters", len(key))
			return
		}
	}

	// validate input length to avoid excessive processing or abuse
	if inputExceedsLimit(rawInput, MaxInputLength) {
		respondWithError(w, http.StatusRequestEntityTooLarge, formatStatusMessage(http.StatusRequestEntityTooLarge, MsgInputTooLong), &data)
		log.Printf("cypherHandler: Input too long: %d characters", len(rawInput))
		return
	}

	// prepopulate the form fields for response rendering
	data.Input = rawInput
	data.Key = key
	
	var result string

	// process input depending on mode
	switch mode {
	case xor:
		// perform XOR encryption/decryption with the provided key
		res, err := functions.Xorify(rawInput, key)
		if err != nil {
			log.Printf("cypherHandler: XOR error: %v", err)
			respondWithError(w, http.StatusInternalServerError, MsgInternalServerError, &data)
			return
		}
		result = res
		data.Mode = xor
		// save successful operation to history
		saveCypherHistory(xor, key, rawInput, result)

	case rot13:
		// perform ROT13 encryption/decryption
		result = functions.Rot13ify(rawInput)
		data.Mode = rot13
		// save successfull operation to history
		saveCypherHistory(rot13, "", rawInput, result)
		
	default:
		// invalid mode, respond with error
		log.Printf("cypherHandler: Invalid mode received: %s", mode)
		respondWithError(w, http.StatusBadRequest, MsgInvalidAction, &data)
		return
	}

	// prepare success response: clear input field, display result
	data.Input = ""
	data.Result = result
	data.StatusCode = http.StatusOK
	data.StatusType = statusSuccess
	data.StatusMessage = formatStatusMessage(http.StatusOK, "successfully encrypted/decrypted")
	data.LineCount = countLines(result)

	// copy the history slice under mutex lock to safely pass to template
	cypherHistoryMutex.Lock()
	data.CypherHistory = make([]CypherHistoryEntry, len(cypherHistory))
	copy(data.CypherHistory, cypherHistory)
	cypherHistoryMutex.Unlock()

	// render the template with updated data and history.
	renderTemplate(w, data)
	
}

/* 
	saveCypherHistory appends a new cypher operation record to the history buffer.
	- uses a mutex to avoid concurrent access issues.
	- keeps the newest entries ath the front of the slice.
	- truncates the history to the last MaxHistoryEntries to limit memory usage.
*/
func saveCypherHistory(mode, key, input, result string) {
	entry := CypherHistoryEntry {
		Timestamp: 	time.Now().Format("January 2, 15:04"),
		Mode:		mode,
		Key:		key,
		Input:		input,
		Result:		result,
	}

	cypherHistoryMutex.Lock()
	defer cypherHistoryMutex.Unlock()

	// insert new entry at the beginning of the slice.
	cypherHistory = append([]CypherHistoryEntry{entry}, cypherHistory...)

	// limit history size to last MaxHistoryEntries entries.
	if len(cypherHistory) > MaxHistoryEntries {
		cypherHistory = cypherHistory[:MaxHistoryEntries]
	}
}