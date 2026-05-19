package main

import (
	"encoding/json"
	"net/http"
)

type AppResponse struct {
	Status int
	Body   any
}

func newAppResponse(status int, body any) *AppResponse {
	return &AppResponse{Status: status, Body: body}
}

func writeJSON(w http.ResponseWriter, status int, requestID string, v any) {
	_, err := json.Marshal(v)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorMessageJSON{
			ErrorCode: INTERNAL_SERVER_ERROR,
			Message:   "Json Decode failed",
			RequestID: requestID,
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}
