package main

import (
	"errors"
	"net/http"
)

type ErrorMessageJSON struct {
	ErrorCode string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id"`
}

func writeError(w http.ResponseWriter, status int, e ErrorMessageJSON) {
	writeJSON(w, status, map[string]ErrorMessageJSON{"error": e})
}

var ErrIncidentNotFound = errors.New("Incident not found")
var ErrIncidentConflict = errors.New("The Incident is already resolved")
var ErrIncidentEmptyList = errors.New("No Incident found")

// Will use if we have database
var ErrInternal = errors.New("Internal Error")

const (
	INCIDENT_NOT_FOUND    = "INCIDENT_NOT_FOUND"
	INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	BAD_REQUEST           = "BAD_REQUEST"
	CONFLICT              = "CONFLICT"
)
