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

func prepareJSON(w http.ResponseWriter, statusCode uint16) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(int(statusCode))
}

func formatJSON(data interface{}, message string, statusCode uint16) ResponseJSON {
	return ResponseJSON{Message: message, Data: data, StatusCode: statusCode}
}

func returnJSON(w http.ResponseWriter, data interface{}) {
	notok := json.NewEncoder(w).Encode(data)
	if notok != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{ "message": "INTERNAL ERROR", "status": 500, data: { "error": "Everything went wrong" } }`)
		return
	}
}
