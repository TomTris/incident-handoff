package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func getRouter(incHandler *IncidentHandler, mongoClient *mongo.Client, promRegistry *prometheus.Registry) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /incidents", ResponseMiddleware(incHandler.CreateIncident))
	mux.HandleFunc("POST /incidents/{id}/entries", ResponseMiddleware(incHandler.AddEntry))
	mux.HandleFunc("GET /incidents/{id}", ResponseMiddleware(incHandler.GetIncident))
	mux.HandleFunc("GET /incidents", ResponseMiddleware(incHandler.ListIncidents))
	mux.HandleFunc("GET /incidents/{id}/handoff", ResponseMiddleware(incHandler.GetHandoffBrief))
	mux.HandleFunc("PATCH /incidents/{id}", ResponseMiddleware(incHandler.UpdateIncident))
	mux.HandleFunc("GET /incidents/{id}/ws", incHandler.HandleIncidentWebSocket)

	mux.HandleFunc("GET /healthz", healthCheck)
	mux.HandleFunc("GET /readyz", readyCheck(mongoClient))
	mux.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{Registry: promRegistry}))

	mux.HandleFunc("POST /flags", incHandler.CreateFlag)
	mux.HandleFunc("GET /flags", incHandler.ListAllFlag)
	mux.HandleFunc("PATCH /flags/{name}", incHandler.UpdateFlag)
	mux.HandleFunc("GET /flags/{name}/evaluate", incHandler.Evaluate)
	router := RequestIDMiddleware(ObservabilityMiddleware(CORSMiddleware(TimeoutMiddleware(mux))))
	return router
}
