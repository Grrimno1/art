package server

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