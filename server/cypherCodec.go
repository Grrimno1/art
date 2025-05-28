package server

import (
	"net/http"
	"art/functions"
	"log"
	"sync"
	"time"
)
// struct to keep Cypher's history of XOR and ROT13 operations
type CypherHistoryEntry struct {
	Timestamp 	string
	Mode		string
	Key			string //optional: used for XOR
	Input		string
	Result		string
}
const (
	xor		= "xor"
	rot13	= "rot13"
)
var (
	cypherHistory 		[]CypherHistoryEntry
	cypherHistoryMutex	sync.Mutex
)

/* handles post requests for XOR and ROT13
	must be x-www-form-urlencoded and contain
		'mode' - determines whether to XOR or ROT13 the input
		'key' - for XOR
		'input' - data to encrypt/decrypt
*/
func CypherHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		log.Printf("Method Not Allowed: received %s, only POST allowed", r.Method)
		return
	}

	data := CombinedPageData{
		Section: "cypher",
	}

	mode := r.FormValue("mode")
	key := r.FormValue("key")
	rawInput := normalizeNewLines(r.FormValue("input"))

	//input validation
	if rawInput == "" {
		http.Error(w, MsgInputEmpty, http.StatusBadRequest)
		log.Printf("Empty input received")
		return
	}
	if mode == xor {
		//key validation not empty and not too long (256)
		if key == "" {
			http.Error(w, MsgXOREmpty, http.StatusBadRequest)
			log.Printf("XOR mode selected but empty key provided")
		}
		if inputExceedsLimit(key, maxKeyLength) {
			data.StatusCode = http.StatusRequestEntityTooLarge
			data.StatusType = StatusError
			data.StatusMessage = formatStatusMessage(http.StatusRequestEntityTooLarge, "XOR key is too long")
			log.Printf("XOR key too long: %d characters", len(key))
			renderTemplate(w, data)
			return
		}
	}
	if inputExceedsLimit(rawInput, MaxInputLength) {
		data.StatusCode = http.StatusRequestEntityTooLarge
		data.StatusType = StatusError
		data.StatusMessage = formatStatusMessage(http.StatusRequestEntityTooLarge, MsgInputTooLong)
		log.Printf("Input too long: %d characters", len(rawInput))
		renderTemplate(w, data)
		return
	}
	//prepopulate for rendering
	data.Input = rawInput
	data.Key = key
	
	var result string

	switch mode {
	case xor:
		res, err := functions.Xorify(rawInput, key)
		if err != nil {
			http.Error(w, "XOR error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("XOR error: %v", err)
			return
		}
		result = res
		data.Mode = xor
		saveCypherHistory(xor, key, rawInput, result)

	case rot13:
		result = functions.Rot13ify(rawInput)
		data.Mode = rot13
		saveCypherHistory(rot13, "", rawInput, result)
		
	default:
		http.Error(w, "Invalid mode: Use 'xor' or 'rot13'", http.StatusBadRequest)
		return
	}
	// populate success response
	data.Result = result
	data.StatusCode = http.StatusOK
	data.StatusType = statusSuccess
	data.StatusMessage = formatStatusMessage(http.StatusOK, "successfully encrypted/decrypted")
	data.LineCount = countLines(result)

	//copying history safely.
	cypherHistoryMutex.Lock()
	data.CypherHistory = make([]CypherHistoryEntry, len(cypherHistory))
	copy(data.CypherHistory, cypherHistory)
	cypherHistoryMutex.Unlock()

	renderTemplate(w, data)
	
}
// saves entry to cypher history buffer.
// keep only the latest MaxHistoryEntries
func saveCypherHistory(mode, key, input, result string) {
	entry := CypherHistoryEntry {
		Timestamp:	time.Now().Format("2006-01-02 15:04"),
		Mode:		mode,
		Key:		key,
		Input:		input,
		Result:		result,
	}
	cypherHistoryMutex.Lock()
	defer cypherHistoryMutex.Unlock()

	//append the new entry
	cypherHistory = append(cypherHistory, entry)
	
	// keep track of the last 20 entries.
	if len(cypherHistory) > MaxHistoryEntries {
		cypherHistory = cypherHistory[len(cypherHistory)-MaxHistoryEntries:]
	}
}