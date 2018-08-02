package global

import (
	"net/http"
)

const (
	NOT_FOUND   = http.StatusNotFound
	BAD_REQUEST = http.StatusBadRequest
)

type ErrorMsg struct {
	Error errorMsg `json:"error"`
}

type errorMsg struct {
	StatusCode int         `json:"statusCode"`
	Name       string      `json:"name"`
	Message    interface{} `json:"message"`
}

func NewErrorMsg(statusCode int, name string, message interface{}) (err *ErrorMsg) {
	return &ErrorMsg{Error: errorMsg{StatusCode: statusCode, Name: name, Message: message}}
}