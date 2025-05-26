package server

import (
	"net/http"
	"html/template"
	"art/functions"
	"strings"
	"sync"
	"time"
)
//own struct to keep Cypher's history separate from artCodec.
type CypherHistoryEntry struct {
	Timestamp 	string
	Mode		string
	Key			string //optional: used for XOR
	Input		string
	Result		string
}

var (
	cypherHistory 		[]CypherHistoryEntry
	cypherHistoryMutex	sync.Mutex
)

var cypherTmpl = template.Must(template.ParseFiles("public/index.html"))
/* handles post requests for XOR and ROT13
	must be x-www-form-urlencoded and contain
		'mode' - determines whether to XOR or ROT13 the input
		'key' - for XOR
		'input' - data to encrypt/decrypt
*/
func CypherHandler(w http.ResponseWriter, r *http.Request) {
	data := CombinedPageData{}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	mode := r.FormValue("mode")
	key := r.FormValue("key")
	input := strings.ReplaceAll(r.FormValue("input"), "\r\n", "\n")
	data.Section = "cypher"

	if input == "" {
		http.Error(w, "Input cannot be empty", http.StatusBadRequest)
		return
	}
	if mode == "xor" && key == "" {
		http.Error(w, "XOR mode requires a non-empty key", http.StatusBadRequest)
		return
	}

	data.Input = input
	data.Key = key
	
	var result string

	switch mode {
	case "xor":
		res, err := functions.Xorify(input, key)
		if err != nil {
			http.Error(w, "XOR error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		result = res
		data.Mode = "xor"
		saveCypherHistory("xor", key, input, result)
	case "rot13":
		result = functions.Rot13ify(input)
		data.Mode = "rot13"
		saveCypherHistory("rot13", "", input, result)
		
	default:
		http.Error(w, "Invalid mode: Use 'xor' or 'rot13'", http.StatusBadRequest)
		return
	}

	data.Result = result
	data.StatusCode = http.StatusOK
	data.StatusType = "success"
	data.StatusMessage = "200 OK: Operation successful"
	data.LineCount = CountLines(result)

	//copying history safely.
	cypherHistoryMutex.Lock()
	data.CypherHistory = make([]CypherHistoryEntry, len(cypherHistory))
	copy(data.CypherHistory, cypherHistory)
	cypherHistoryMutex.Unlock()

	if err := cypherTmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
	
}

func saveCypherHistory(mode, key, input, result string) {
	entry := CypherHistoryEntry {
		Timestamp:	time.Now().Format("02.01.06 15:04"),
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
	if len(cypherHistory) > 20 {
		cypherHistory = cypherHistory[len(cypherHistory)-20:]
	}
}