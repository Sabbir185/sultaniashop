package httpx

import (
	"encoding/json"
	"net/http"
)

type StatusCode string

const (
	CodeInternalError   StatusCode = "internal_error"
	CodeNotFound        StatusCode = "not_found"
	CodeBadRequest      StatusCode = "bad_request"
	CodeForbidden       StatusCode = "forbidden"
	CodeUnauthorized    StatusCode = "unauthorized"
	CodeConflict        StatusCode = "conflict"
	CodeTooManyRequests StatusCode = "too_many_requests"
)

type errorPayload struct {
	Code    StatusCode `json:"code"`
	Message string     `json:"message"`
}

type errorEnvelope struct {
	Error errorPayload `json:"error"`
}

func Error(w http.ResponseWriter, status int, message string, code StatusCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{
		Error: errorPayload{
			Code:    code,
			Message: message,
		},
	})
}
