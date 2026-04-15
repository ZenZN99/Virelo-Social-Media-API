package app

import (
	"backend/utils"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Container struct {
	Handlers    *Handlers
	Middlewares *Middlewares
	Services    *Services
}

func NewContainer() (*Container, error) {

	mongoURI := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return &Container{}, err
	}

	dbName := os.Getenv("DB_NAME")
	db := client.Database(dbName)
	userCollection := db.Collection("users")
	postCollection := db.Collection("posts")
	reelCollection := db.Collection("reels")
	commentCollection := db.Collection("comments")
	likeCollection := db.Collection("likes")
	followCollection := db.Collection("follows")
	messageCollection := db.Collection("messages")
	notificationCollection := db.Collection("notifications")

	tokenService := utils.NewTokenService()

	// Init Services
	services := InitServices(userCollection, postCollection, reelCollection, commentCollection, likeCollection, followCollection, messageCollection, notificationCollection)

	// Init Handlers
	handlers := InitHandlers(services)
	// Init Middlewares
	middlewares := InitMiddlewares(tokenService)

	return &Container{
		Handlers:    handlers,
		Middlewares: middlewares,
		Services:    services,
	}, nil
}
