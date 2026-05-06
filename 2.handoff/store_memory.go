package main

import (
	"context"
	"errors"
	"strconv"
	"time"
)

type MemoryStore struct {
	incidents           map[string]Incident
	nextIncidentID      int
	nextEntryTimelineID int
}

func (m *MemoryStore) CreateIncident(ctx context.Context, inc Incident) (Incident, error) {
	m.nextIncidentID++
	inc.ID = incidentIDPrefix + strconv.Itoa(m.nextIncidentID)
	inc.Status = TRIGGERED
	inc.CreatedAt = time.Now()
	inc.UpdatedAt = time.Now()

	m.incidents[inc.ID] = inc
	return inc, nil
}

func (m *MemoryStore) GetIncident(ctx context.Context, id string) (Incident, error) {
	inc, ok := m.incidents[id]
	if ok == false {
		return inc, errors.New("Incident not found")
	}
	return inc, nil
}

// func defaultNewTimelineEntry() TimelineEntry {
// 	nextEntryTimelineID++
// 	return TimelineEntry{
// 		ID:   entryIDPrefix + strconv.Itoa(nextEntryTimelineID),
// 		Time: time.Now(),
// 	}
// }

// func addTimelineEntry(w http.ResponseWriter, r *http.Request) {
// 	newTimelineEntry := defaultNewTimelineEntry()
// 	err := json.NewDecoder(r.Body).Decode(&newTimelineEntry)
// 	if err != nil {
// 		http.Error(w, "failed to read body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	incidentId := r.PathValue("id")
// 	for idx, _ := range incidents {
// 		if incidents[idx].ID == incidentId {
// 			incidents[idx].Entries = append(incidents[idx].Entries, newTimelineEntry)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.WriteHeader(http.StatusCreated)
// 			json.NewEncoder(w).Encode(newTimelineEntry)
// 			return
// 		}
// 	}
// 	http.Error(w, "incident not found", http.StatusNotFound)
// }

// func listAllIncidents(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "allication/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(incidents)
// }
