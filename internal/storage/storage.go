package storage

import (
	"context"
	"log"
	"time"

	"100.GO/internal/http/handler"
	"100.GO/internal/infrastructura/repository/mongodb"
	"100.GO/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewStorage() (*mongo.Client, *mongo.Collection, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	collection := client.Database("N9").Collection("n9")
	return client, collection, nil
}

func Handler() *handler.UserHandler {
	client, collection, err := NewStorage()
	if err != nil {
		log.Println("connection mongodb error")
	}
	repo := mongodb.NewUserMongodb(client,collection)

	service := service.NewUserService(repo)

	handler := handler.NewUserHandler(service)

	return handler
}
