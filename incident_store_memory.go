package main

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"
)

type MemoryIncidentStore struct {
	mu                  sync.RWMutex
	incidents           map[string]Incident
	nextIncidentID      int
	nextEntryTimelineID int
}

func NewMemoryIncidentStore() *MemoryIncidentStore {
	return &MemoryIncidentStore{incidents: make(map[string]Incident)}
}

func (m *MemoryIncidentStore) CreateIncident(ctx context.Context, openedBy string, onCall string, incReq CreateIncidentRequest) (Incident, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nextIncidentID++
	inc := Incident{
		ID:        incidentIDPrefix + strconv.Itoa(m.nextIncidentID),
		Title:     incReq.Title,
		Service:   incReq.Service,
		Severity:  incReq.Severity,
		OpenedBy:  openedBy,
		OnCall:    onCall,
		Status:    TRIGGERED,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Entries:   []TimelineEntry{},
		Version:   1,
	}
	m.incidents[inc.ID] = inc
	return inc, nil
}

func (m *MemoryIncidentStore) GetIncident(ctx context.Context, id string) (Incident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	inc, ok := m.incidents[id]
	if ok == false {
		return inc, ErrIncidentNotFound
	}
	return inc, nil
}

func (m *MemoryIncidentStore) AddEntry(ctx context.Context, incID string, currentIncVersion int, entry TimelineEntry) (TimelineEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	inc, ok := m.incidents[incID]
	if ok == false {
		return TimelineEntry{}, ErrIncidentNotFound
	}
	if inc.Status == RESOLVED {
		return TimelineEntry{}, ErrIncidentResolved
	}
	if inc.Version != currentIncVersion {
		return TimelineEntry{}, ErrIncidentVersionConflict
	}

	m.nextEntryTimelineID++
	entry.ID = entryIDPrefix + strconv.Itoa(m.nextEntryTimelineID)
	entry.CreatedAt = time.Now()
	inc.Entries = append(inc.Entries, entry)
	inc.UpdatedAt = time.Now()
	inc.Version = currentIncVersion + 1
	m.incidents[incID] = inc
	return entry, nil
}

func (m *MemoryIncidentStore) ListIncidents(ctx context.Context, filter IncidentFilter) ([]Incident, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	isServiceMatch := func(incident Incident, filter IncidentFilter) bool {
		return filter.Service == "" || filter.Service == incident.Service
	}
	isStatusMatch := func(incident Incident, filter IncidentFilter) bool {
		return filter.Status == "" ||
			(filter.Status == "active" && incident.Status != RESOLVED) ||
			filter.Status == incident.Status
	}

	array := []Incident{}
	for _, incident := range m.incidents {
		if isServiceMatch(incident, filter) && isStatusMatch(incident, filter) {
			array = append(array, incident)
		}
	}
	sort.Slice(array, func(i, j int) bool {
		return array[i].CreatedAt.Sub(array[j].CreatedAt) < 0
	})
	return array, nil
}

// Return incident After
func (m *MemoryIncidentStore) UpdateIncident(ctx context.Context, incID string, currentIncVersion int, update IncidentUpdate) (Incident, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	incident, ok := m.incidents[incID]
	if ok == false {
		return incident, ErrIncidentNotFound
	}
	if incident.Version != currentIncVersion {
		return incident, ErrIncidentVersionConflict
	}

	if update.Status != nil {
		incident.Status = *update.Status
	}
	if update.Severity != nil {
		incident.Severity = *update.Severity
	}
	if update.OnCall != nil {
		incident.OnCall = *update.OnCall
	}

	incident.Version = currentIncVersion + 1
	incident.UpdatedAt = time.Now()
	m.incidents[incID] = incident
	return incident, nil
}
