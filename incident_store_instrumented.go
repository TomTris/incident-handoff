package main

import (
	"context"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type InstrumentedIncidentStore struct {
	inner   IncidentStore
	metrics *IncidentStoreMetrics
}

func (s *InstrumentedIncidentStore) MetricInit() {
	incidents, err := s.inner.ListIncidents(context.Background(), IncidentFilter{})
	if err != nil {
		log.Fatal(err)
	}
	for _, incident := range incidents {
		s.metrics.IncidentTotal.WithLabelValues(incident.Status).Inc()
		s.metrics.TotalEntries.Add(float64(len(incident.Entries)))
	}
}

func (s *InstrumentedIncidentStore) CreateIncident(ctx context.Context, openedBy string, onCall string, incReq CreateIncidentRequest) (Incident, error) {
	timer := prometheus.NewTimer(s.metrics.DbQueryDurationSeconds.WithLabelValues("create_incident"))
	inc, err := s.inner.CreateIncident(ctx, openedBy, onCall, incReq)
	timer.ObserveDuration()
	if err == nil {
		s.metrics.IncidentTotal.WithLabelValues(inc.Status).Inc()
	}
	return inc, err
}

func (s *InstrumentedIncidentStore) GetIncident(ctx context.Context, id string) (Incident, error) {
	timer := prometheus.NewTimer(s.metrics.DbQueryDurationSeconds.WithLabelValues("get_incident"))
	defer timer.ObserveDuration()
	return s.inner.GetIncident(ctx, id)
}

func (s *InstrumentedIncidentStore) AddEntry(ctx context.Context, incID string, expectedIncVersion int, entry TimelineEntry) (TimelineEntry, error) {
	timer := prometheus.NewTimer(s.metrics.DbQueryDurationSeconds.WithLabelValues("add_entry"))
	entry, err := s.inner.AddEntry(ctx, incID, expectedIncVersion, entry)
	timer.ObserveDuration()

	if err == nil {
		s.metrics.TotalEntries.Inc()
	}
	return entry, err
}

func (s *InstrumentedIncidentStore) ListIncidents(ctx context.Context, filter IncidentFilter) ([]Incident, error) {
	timer := prometheus.NewTimer(s.metrics.DbQueryDurationSeconds.WithLabelValues("list_incident"))
	defer timer.ObserveDuration()
	return s.inner.ListIncidents(ctx, filter)
}

func (s *InstrumentedIncidentStore) UpdateIncident(ctx context.Context, incID string, expectedIncVersion int, update IncidentUpdate) (Incident, error) {
	timer := prometheus.NewTimer(s.metrics.DbQueryDurationSeconds.WithLabelValues("update_incident"))
	incAfter, err := s.inner.UpdateIncident(ctx, incID, expectedIncVersion, update)
	timer.ObserveDuration()
	if err == nil && update.Status != nil {
		s.metrics.IncidentTotal.WithLabelValues(*update.Status).Inc()
		s.metrics.IncidentTotal.WithLabelValues(incAfter.Status).Dec()
	}
	return incAfter, err
}
