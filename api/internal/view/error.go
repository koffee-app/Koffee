package view

import (
	"koffee/internal/models"
	"net/http"
)

// ErrorResp Response error
type ErrorResp interface {
	Message() string
}

func end(w http.ResponseWriter, err ErrorResp, code uint16) {
	ReturnJSON(w, formatJSON(err, err.Message(), code))
}

// RenderAuthError Sends to the client the error
func RenderAuthError(w http.ResponseWriter, u *models.UserError) {
	prepareJSON(w, http.StatusBadRequest)
	end(w, u, http.StatusBadRequest)
}

// InternalError Call this when you wanna send an internal error.
func InternalError(w http.ResponseWriter, data interface{}) {
	response := ResponseJSON{StatusCode: http.StatusInternalServerError, Message: "Internal server error.", Data: data}
	ReturnJSON(w, response)
}

// ErrorAuthentication Respond to the client a JSON error
func ErrorAuthentication(w http.ResponseWriter, err interface{}) {
	response := ResponseJSON{StatusCode: http.StatusUnauthorized, Message: "Error, invalid token.", Data: err}
	prepareJSON(w, response.StatusCode)
	ReturnJSON(w, response)
}

// Error Returns an error with no standard message.
func Error(w http.ResponseWriter, message string, code uint16, err interface{}) {
	response := ResponseJSON{StatusCode: code, Message: message, Data: err}
	prepareJSON(w, code)
	ReturnJSON(w, response)
}
