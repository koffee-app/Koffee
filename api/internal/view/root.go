package view

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseJSON Response Struct
type ResponseJSON struct {
	StatusCode uint16      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

// PrepareJSON Adds the headers.
func PrepareJSON(w http.ResponseWriter, statusCode uint16) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(int(statusCode))
}

// FormatJSON Formats the response in the Koffee server standard
func FormatJSON(data interface{}, message string, statusCode uint16) ResponseJSON {
	return ResponseJSON{Message: message, Data: data, StatusCode: statusCode}
}

// ReturnJSON Returns JSON
func ReturnJSON(w http.ResponseWriter, data interface{}) {
	notok := json.NewEncoder(w).Encode(data)
	if notok != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "message": "INTERNAL ERROR", "status": 500, data: { "error": "Everything went wrong" } }`)
		return
	}
}
