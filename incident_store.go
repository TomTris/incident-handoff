package main

import "context"

type IncidentStore interface {
	CreateIncident(ctx context.Context, onCall string, incReq CreateIncidentRequest) (Incident, error)
	GetIncident(ctx context.Context, id string) (Incident, error)
	ListIncidents(ctx context.Context, filter IncidentFilter) ([]Incident, error)
	UpdateIncident(ctx context.Context, incID string, expectedIncVersion int, update IncidentUpdate) (Incident, error) // Return incident After
	AddEntry(ctx context.Context, incID string, expectedIncVersion int, entry TimelineEntry) (TimelineEntry, error)
}
