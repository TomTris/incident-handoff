package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config := loadConfig()
	var store Store
	var storeType string
	if config.ConnectionString != "" {
		store = NewMongoStore(config.ConnectionString, config.DatabaseName)
		storeType = "MongoStore"
	} else {
		store = &MemoryStore{incidents: make(map[string]Incident)}
		slog.Info("no connection string, using in-memory store")
		storeType = "MemoryStore"
	}
	incHandler := IncidentHandler{Store: store}
	router := getRouter(incHandler)

	srv := http.Server{
		Addr:    ":" + config.Port,
		Handler: router,
	}

	go func() {
		slog.Info(fmt.Sprintf("server starting port=%s, store=%s", srv.Addr, storeType))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	slog.Info("server shut down in <= 10 sec")
	srv.Shutdown(ctx)
	slog.Info("server shut down gracefully")
}
