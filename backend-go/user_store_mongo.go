package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserStore struct {
	db *mongo.Database
}

func NewMongoUserStore(client *mongo.Client, DBName string) *MongoUserStore {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}

	db := client.Database(DBName)

	_, err := db.Collection(CollectionUsers).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		log.Fatal("ensure users unique index: " + err.Error())
	}
	slog.Info("user schema/indexes ensured")
	return &MongoUserStore{db: db}
}

func (s *MongoUserStore) nextID(ctx context.Context, name string, prefix string) (string, error) {
	col := s.db.Collection(CollectionUserCounters)
	filter := bson.M{"_id": name}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result struct {
		Seq int `bson:"seq"`
	}
	err := col.FindOneAndUpdate(ctx, filter, update, opts).Decode(&result)
	if err != nil {
		return "", err
	}
	return prefix + strconv.Itoa(result.Seq), nil
}

func (s *MongoUserStore) Create(ctx context.Context, u User) (User, error) {
	id, err := s.nextID(ctx, "user", UserPrefix)
	if err != nil {
		return User{}, fmt.Errorf("next user id: %w", err)
	}
	u.ID = id

	col := s.db.Collection(CollectionUsers)
	_, err = col.InsertOne(ctx, u)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return User{}, ErrUserAlreadyExist
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	return u, nil
}

func (s *MongoUserStore) GetByUsername(ctx context.Context, username string) (User, error) {
	col := s.db.Collection(CollectionUsers)
	var u User
	err := col.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return u, nil
}
