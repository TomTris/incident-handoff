package main

import "net/http"

type ErrorMessageJSON struct {
	ErrorCode string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func writeError(w http.ResponseWriter, status int, e ErrorMessageJSON) {
	writeJSON(w, status, map[string]ErrorMessageJSON{"error": e})
}
