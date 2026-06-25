package main

import (
	"errors"
	"testing"
	"time"
)

func runStoreContractsTests(t *testing.T, makeStore func(t *testing.T) IncidentStore) {
	t.Run("CreateIncident", func(t *testing.T) { TestIncidentStoreCreateIncident(t, makeStore) })
	t.Run("UpdateIncident", func(t *testing.T) { TestIncidentStoreUpdateIncident(t, makeStore) })
	t.Run("AddEntry", func(t *testing.T) { TestIncidentStoreAddEntry(t, makeStore) })
	t.Run("ListIncidents", func(t *testing.T) { TestIncidentStoreListIncidents(t, makeStore) })
}

func TestIncidentStoreCreateIncident(t *testing.T, makeStore func(t *testing.T) IncidentStore) {
	m := makeStore(t)

	t.Run("normal creation", func(t *testing.T) {
		inc, err := m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
			Title:    "outage",
			Service:  "api",
			Severity: "SEV1",
		})
		if err != nil {
			t.Fatal(err)
		}
		if inc.Title != "outage" {
			t.Errorf("Title expected `outage`, got `%s`", inc.Title)
		}
	})

	t.Run("sets correct defaults", func(t *testing.T) {
		inc, _ := m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
			Title:    "outage3",
			Service:  "api",
			Severity: "SEV1",
		})
		if inc.Status != TRIGGERED {
			t.Errorf("Status expected %s, got %s", TRIGGERED, inc.Status)
		}
		if inc.CreatedAt.IsZero() {
			t.Errorf("CreateAt not set")
		}
		if len(inc.Entries) != 0 {
			t.Errorf("len(inc.Entries) expected 0, got %v", len(inc.Entries))
		}
		if inc.Version != 1 {
			t.Errorf("Version expected 1, got %v", inc.ID)
		}
	})

	t.Run("map incident fields correcy", func(t *testing.T) {
		inc, _ := m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
			Title:    "outage4",
			Service:  "api",
			Severity: "SEV1",
		})
		if inc.Title != "outage4" {
			t.Errorf("Title expected %s, got %s", "outage4", inc.Title)
		}
		if inc.Service != "api" {
			t.Errorf("Service expected %s, got %s", "api", inc.Service)
		}
		if inc.Severity != "SEV1" {
			t.Errorf("Severity expected %s, got %s", "SEV", inc.Severity)
		}
	})
	t.Run("set openedBy and OnCall correctly", func(t *testing.T) {
		inc, _ := m.CreateIncident(t.Context(), "openByAnh", "OnCallMega", CreateIncidentRequest{
			Title:    "outage4",
			Service:  "api",
			Severity: "SEV1",
		})
		if inc.OpenedBy != "openByAnh" {
			t.Errorf("OpenedBy expected %s, got %s", "openByAnh", inc.OpenedBy)
		}
		if inc.OnCall != "OnCallMega" {
			t.Errorf("Service expected %s, got %s", "OnCallMega", inc.OnCall)
		}
	})
	t.Run("sequential IDs", func(t *testing.T) {
		m = makeStore(t)
		inc1, _ := m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
			Title: "a", Service: "s", Severity: "SEV1",
		})
		inc2, _ := m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
			Title: "b", Service: "s", Severity: "SEV1",
		})
		if inc1.ID != "INC-1" {
			t.Errorf("expected %s, got %s", "INC-1", inc1.ID)
		}
		if inc2.ID != "INC-2" {
			t.Errorf("expected %s, got %s", "INC-2", inc2.ID)
		}
	})
}

func TestIncidentStoreUpdateIncident(t *testing.T, makeStore func(t *testing.T) IncidentStore) {
	m := makeStore(t)

	m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
		Title: "outage", Service: "api", Severity: "SEV1",
	})

	t.Run("update status", func(t *testing.T) {
		inc, err := m.UpdateIncident(t.Context(), "INC-1", 1, IncidentUpdate{
			Status: new(RESOLVED),
		})
		if err != nil {
			t.Fatal(err)
		}
		if inc.Status != RESOLVED {
			t.Errorf("expected %s, got %s", RESOLVED, inc.Status)
		}
	})

	t.Run("test sequential Version", func(t *testing.T) {
		_, err := m.UpdateIncident(t.Context(), "INC-1", 1, IncidentUpdate{
			Status: new(MITIGATED),
		})
		if err != ErrIncidentVersionConflict {
			t.Fatalf("expect Version conflict, get %v", err)
		}
	})

	t.Run("update severity", func(t *testing.T) {
		inc, err := m.UpdateIncident(t.Context(), "INC-1", 2, IncidentUpdate{
			Severity: new("SEV2"),
		})
		if err != nil {
			t.Fatal(err)
		}
		if inc.Severity != "SEV2" {
			t.Errorf("expected SEV2, got %s", inc.Severity)
		}
	})

	t.Run("update on_call", func(t *testing.T) {
		inc, err := m.UpdateIncident(t.Context(), "INC-1", 3, IncidentUpdate{
			OnCall: new("tom"),
		})
		if err != nil {
			t.Fatal(err)
		}
		if inc.OnCall != "tom" {
			t.Errorf("expected tom, got %s", inc.OnCall)
		}
	})

	t.Run("multiple fields at once", func(t *testing.T) {
		inc, err := m.UpdateIncident(t.Context(), "INC-1", 4, IncidentUpdate{
			Status:   new(TRIGGERED),
			Severity: new("SEV3"),
			OnCall:   new("carl"),
		})
		if err != nil {
			t.Fatal(err)
		}
		if inc.Status != TRIGGERED {
			t.Errorf("Status expected %s, got %s", TRIGGERED, inc.Status)
		}
		if inc.Severity != "SEV3" {
			t.Errorf("Severity expected SEV3, got %s", inc.Severity)
		}
		if inc.OnCall != "carl" {
			t.Errorf("OnCall expected carl, got %s", inc.OnCall)
		}
	})

	t.Run("updated_at changes", func(t *testing.T) {
		before, _ := m.GetIncident(t.Context(), "INC-1")
		time.Sleep(10 * time.Millisecond) //Give Buffer to compare updated_at
		after, _ := m.UpdateIncident(t.Context(), "INC-1", 5, IncidentUpdate{
			Status: new(TRIGGERED),
		})
		if after.UpdatedAt.After(before.UpdatedAt) == false {
			t.Error("UpdatedAt should advance")
		}
	})
}

func TestIncidentStoreAddEntry(t *testing.T, makeStore func(t *testing.T) IncidentStore) {
	m := makeStore(t)

	m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
		Title:    "outage",
		Service:  "api",
		Severity: "SEV1",
	})

	t.Run("adds entry to existing incident", func(t *testing.T) {
		entry, err := m.AddEntry(t.Context(), "INC-1", 1, TimelineEntry{
			Author: "anh",
			Type:   "observation",
			Text:   "pool exhausted",
		})
		if err != nil {
			t.Fatal(err)
		}
		if entry.ID != "TLE-1" {
			t.Errorf("expected %s, got %s", "TLE-1", entry.ID)
		}
		if entry.CreatedAt.IsZero() {
			t.Error("Time not set")
		}
	})
	t.Run("sequential entry IDs", func(t *testing.T) {
		entry, _ := m.AddEntry(t.Context(), "INC-1", 2, TimelineEntry{
			Author: "anh", Type: "observation", Text: "second entry",
		})
		if entry.ID != "TLE-2" {
			t.Errorf("expected %s, got %s", "TLE-2", entry.ID)
		}
	})

	t.Run("conflict on resolved incident", func(t *testing.T) {
		// resolve the incident first
		m.UpdateIncident(t.Context(), "INC-1", 3, IncidentUpdate{
			Status: new(RESOLVED),
		})

		_, err := m.AddEntry(t.Context(), "INC-1", 4, TimelineEntry{
			Author: "anh", Type: "observation", Text: "too late",
		})
		if !errors.Is(err, ErrIncidentResolved) {
			t.Errorf("expected ErrIncidentConflict, got %v", err.Error())
		}
	})
}

func TestIncidentStoreListIncidents(t *testing.T, makeStore func(t *testing.T) IncidentStore) {
	m := makeStore(t)

	m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
		Title: "a", Service: "api", Severity: "SEV1",
	})
	m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
		Title: "b", Service: "chatbot", Severity: "SEV2",
	})
	m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
		Title: "c", Service: "biling", Severity: "SEV1",
	})
	m.CreateIncident(t.Context(), "", "", CreateIncidentRequest{
		Title: "d", Service: "api", Severity: "SEV3",
	})

	m.UpdateIncident(t.Context(), "INC-1", 1, IncidentUpdate{Status: new(RESOLVED)})
	m.UpdateIncident(t.Context(), "INC-2", 1, IncidentUpdate{Status: new(RESOLVED)})
	m.UpdateIncident(t.Context(), "INC-3", 1, IncidentUpdate{Status: new(INVESTIGATING)})

	t.Run("no filter returns all", func(t *testing.T) {
		list, err := m.ListIncidents(t.Context(), IncidentFilter{})
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 4 {
			t.Errorf("expected 4, got %d", len(list))
		}
	})

	t.Run("filter by service", func(t *testing.T) {
		list, err := m.ListIncidents(t.Context(), IncidentFilter{
			Service: "api",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 2 {
			t.Errorf("expected 2, got %d", len(list))
		}
		for _, inc := range list {
			if inc.Service != "api" {
				t.Errorf("expected service api, got %s", inc.Service)
			}
		}
	})

	t.Run("status active excludes resolved", func(t *testing.T) {
		list, err := m.ListIncidents(t.Context(), IncidentFilter{
			Status: "active",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 2 {
			t.Errorf("expected 2, got %d", len(list))
		}
		for _, inc := range list {
			if inc.Status == RESOLVED {
				t.Errorf("should not include resolved incidents")
			}
		}
	})

	t.Run("specific status", func(t *testing.T) {
		list, err := m.ListIncidents(t.Context(), IncidentFilter{
			Status: RESOLVED,
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 2 {
			t.Errorf("expected 1, got %d", len(list))
		}
		if list[0].ID != "INC-1" {
			t.Errorf("expected INC-1, got %s", list[0].ID)
		}
		if list[1].ID != "INC-2" {
			t.Errorf("expected INC-2, got %s", list[1].ID)
		}
	})

	t.Run("combined service and status", func(t *testing.T) {
		list, err := m.ListIncidents(t.Context(), IncidentFilter{
			Service: "api",
			Status:  "active",
		})
		if err != nil {
			t.Fatal(err)
		}
		if len(list) != 1 {
			t.Errorf("expected 1, got %d", len(list))
		}
		if list[0].ID != "INC-4" {
			t.Errorf("expected INC-4, got %s", list[0].ID)
		}
	})

	t.Run("sorted by created_at ascending", func(t *testing.T) {
		list, _ := m.ListIncidents(t.Context(), IncidentFilter{})
		for i := 1; i < len(list); i++ {
			if list[i].CreatedAt.Before(list[i-1].CreatedAt) {
				t.Errorf("index %d created before index %d", i, i-1)
			}
		}
	})

	t.Run("no match returns empty slice not nil", func(t *testing.T) {
		list, err := m.ListIncidents(t.Context(), IncidentFilter{
			Service: "nonexistent",
		})
		if err != nil {
			t.Fatal(err)
		}
		if list == nil {
			t.Error("expected empty slice, got nil")
		}
		if len(list) != 0 {
			t.Errorf("expected 0, got %d", len(list))
		}
	})
}
