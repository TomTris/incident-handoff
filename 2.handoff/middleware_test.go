package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(requestIDKey)
		if id == nil {
			t.Fatal("no requestID in context")
		}
		w.WriteHeader(http.StatusOK)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	RequestIDMiddleware(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status expected %v, got %v", http.StatusOK, rec.Code)
	}
	if rec.Header().Get("X-Request-ID") == "" {
		t.Fatalf("Header expected %v, got %v", "X-Request-ID", "empty")
	}
}

func TestCORSMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	CORSMiddleware(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status expected %v, got %v", http.StatusOK, rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("Header expected %v, got %v", "*", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Header().Get("Access-Control-Allow-Method") != "GET, POST, PATCH, DELETE" {
		t.Fatalf("Header expected %v, got %v", "GET, POST, PATCH, DELETE", rec.Header().Get("Access-Control-Allow-Method"))
	}
	if rec.Header().Get("Access-Control-Allow-Headers") != "Content-type, Authorization" {
		t.Fatalf("Header expected %v, got %v", "Content-type, Authorization", rec.Header().Get("Access-Control-Allow-Headers"))
	}
}
