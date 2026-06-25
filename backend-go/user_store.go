package main

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserStore interface {
	Create(ctx context.Context, u User) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
}

func NewUsertStore(client *mongo.Client, conf Config) UserStore {
	if conf.ConnectionString == "" {
		slog.Info("use in-memory store for UserStore")
		return NewMemoryUserStore()
	}
	slog.Info("use MongoStore for UserStore")
	return NewMongoUserStore(client, conf.DatabaseName)
}
