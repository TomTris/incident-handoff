package main

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/healthz", nil)
	ctx := context.WithValue(req.Context(), requestIDKey, "test-health-check-id")
	healthCheck(rec, req.WithContext(ctx))
	if rec.Code != 200 {
		t.Fatalf("Expected Code %v, get %v", 200, rec.Code)
	}
	var body map[string]string
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("expected no error, got err %v", err.Error())
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status %v, got %v", "ok", body["status"])
	}
}

func TestReadyCheck(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/readyz", nil)
	ctx := context.WithValue(req.Context(), requestIDKey, "test-ready-check-id")
	readyCheck(nil)(rec, req.WithContext(ctx))
	if rec.Code != 200 {
		t.Fatalf("Expected Code %v, get %v", 200, rec.Code)
	}
	var body map[string]string
	err := json.NewDecoder(rec.Body).Decode(&body)
	if err != nil {
		t.Fatalf("expected no error, got err %v", err.Error())
	}
	if body["status"] != "ready" {
		t.Fatalf("expected status %v, got %v", "ready", body["status"])
	}
}
