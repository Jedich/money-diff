package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"money-diff/bot/helpers"
)

type Connection struct {
	Client *mongo.Client
	Ctx    context.Context
}

//OpenConnection return a connection of a mongodb driver
func OpenConnection(ctx context.Context) *mongo.Client {
	uri := helpers.Getenv("MONGODB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		log.Fatal(err)
	}

	fmt.Println("connected to mongo")
	return client
}
