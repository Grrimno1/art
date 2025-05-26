package server

import (
	"art/functions"
	"html/template"
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
)

var tmpl = template.Must(template.ParseFiles("public/index.html"))

/*
	handling /decoder POST requests for both encoding and decoding
*/

func CodecHandler(w http.ResponseWriter, r *http.Request) {
	data := CombinedPageData{
		Section: "decoder",
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.StatusCode = http.StatusBadRequest
			data.StatusType = "error"
			data.StatusMessage = "400 Bad Request: Failed to parse form"
			w.WriteHeader(http.StatusBadRequest)
			_ = tmpl.Execute(w, data)
			return
		}

		data.DecodeInput = strings.ReplaceAll(r.FormValue("decodeInput"), "\r\n", "\n")
		data.EncodeInput = strings.ReplaceAll(r.FormValue("encodeInput"), "\r\n", "\n")
		action := r.FormValue("action")

		switch action {
		case "encode":
			result := functions.EncodeString(data.EncodeInput, false)
			if result == "Error\n" {
				data.StatusCode = http.StatusBadRequest
				data.StatusType = "error"
				data.StatusMessage = "400 Bad Request: Malformed input"
				data.DecodeInput = ""
				w.WriteHeader(http.StatusBadRequest)
			} else {
				data.DecodeInput = result
				data.StatusCode = http.StatusAccepted
				data.StatusType = "success"
				data.StatusMessage = "202 Accepted: Successfully encoded."
				w.WriteHeader(http.StatusAccepted)
				
				//saving input in history.
				saveHistory("encode", data.EncodeInput, result)
			}

		case "decode":
			result := functions.DecodeString(data.DecodeInput, false)
			if result == "Error\n" {
				data.StatusCode = http.StatusBadRequest
				data.StatusType = "error"
				data.StatusMessage = "400 Bad Request: Malformed input"
				data.EncodeInput = ""
				w.WriteHeader(http.StatusBadRequest)
			} else {
				data.EncodeInput = result
				data.StatusCode = http.StatusAccepted
				data.StatusType = "success"
				data.StatusMessage = "202 Accepted: Successfully decoded."
				w.WriteHeader(http.StatusAccepted)

				//saving input in history.
				saveHistory("decode", data.DecodeInput, result)
			}

		default:
			data.StatusCode = http.StatusBadRequest
			data.StatusType = "error"
			data.StatusMessage = "400 Bad Request: Invalid action"
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	decodedLines := CountLines(data.DecodeInput)
	encodedLines := CountLines(data.EncodeInput)
	if decodedLines > encodedLines {
		data.LineCount = decodedLines
	} else {
		data.LineCount = encodedLines
	}
	//safely reading history before assigning template data.
	historyMutex.Lock()
	data.History = make([]HistoryEntry, len(history))
	copy(data.History, history)
	historyMutex.Unlock()

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
	}
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

//counting the lines of input so we can adjust textareas in frontend for better user-experience.
func CountLines(s string) int {
	lines := strings.Count(s, "\n") + 1
	if lines < 4 {
		return 4
	}
	return lines
}
