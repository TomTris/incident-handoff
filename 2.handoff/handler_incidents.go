package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type IncidentHandler struct {
	Store Store
}

func (incHandler *IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	incidentID := r.PathValue("id")
	inc, err := incHandler.Store.GetIncident(r.Context(), incidentID)
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: "INCIDENT_NOT_fOUND",
			Message:   err.Error(),
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}
	writeJSON(w, http.StatusOK, inc)
}

func (incHandler *IncidentHandler) AddEntry(ctx context.Context, incidentID string, entry TimelineEntry) error {
	return nil
}

func (incHandler *IncidentHandler) CreateIncident(w http.ResponseWriter, r *http.Request) {
	newCreateIncidentRequest := CreateIncidentRequest{}
	err := json.NewDecoder(r.Body).Decode(&newCreateIncidentRequest)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: "BAD_REQUEST",
			Message:   fmt.Sprintf("Validation invalid: %s", err),
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}

	err = newCreateIncidentRequest.Validate()
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorMessageJSON{
			ErrorCode: "BAD_REQUEST",
			Message:   fmt.Sprintf("Validation invalid: %s", err),
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}

	createdIncident, err := incHandler.Store.CreateIncident(r.Context(), Incident{
		Title:    newCreateIncidentRequest.Title,
		Service:  newCreateIncidentRequest.Service,
		Severity: newCreateIncidentRequest.Severity,
		OpenedBy: newCreateIncidentRequest.OpenedBy,
	})

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorMessageJSON{
			ErrorCode: "INTERNAL_ERROR",
			Message:   "failed to create incident",
			RequestID: r.Context().Value(requestIDKey).(string),
		})
		return
	}

	writeJSON(w, http.StatusCreated, createdIncident)
	return
}
