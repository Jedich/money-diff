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

//OpenConnection return a connection of a mongodb driver
func OpenConnection() *mongo.Client {
	uri := helpers.Getenv("MONGODB_URI")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("failed to connect to database: %s", err)
	}

	if err := client.Ping(context.Background(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		log.Fatal(err)
	}

	fmt.Println("connected to mongo")
	return client
}

func CloseConnection(c *mongo.Client) {
	if err := c.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("disconnected from dao")
}
