package main

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IncidentStore interface {
	CreateIncident(ctx context.Context, openedBy string, onCall string, incReq CreateIncidentRequest) (Incident, error)
	GetIncident(ctx context.Context, id string) (Incident, error)
	ListIncidents(ctx context.Context, filter IncidentFilter) ([]Incident, error)
	UpdateIncident(ctx context.Context, incID string, expectedIncVersion int, update IncidentUpdate) (Incident, error) // Return incident After
	AddEntry(ctx context.Context, incID string, expectedIncVersion int, entry TimelineEntry) (TimelineEntry, error)
}

func NewIncidentStore(client *mongo.Client, conf Config) IncidentStore {
	if client == nil {
		slog.Info("use in-memory store for IncidentStore")
		return NewMemoryIncidentStore()
	}
	slog.Info("use MongoStore for IncidentStore")
	return NewMongoIncidentStore(client, conf.DatabaseName)
}
