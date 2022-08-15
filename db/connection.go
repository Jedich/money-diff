package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"chat_id", 1}, {"user_id", 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = client.Database("money").Collection("participants").Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func CloseConnection(c *mongo.Client) {
	if err := c.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("disconnected from repository")
}

//WithTransaction run queries with transaction
//return error if aborted
func WithTransaction(client *mongo.Client, toRun func(*mongo.Client) error) error {
	err := client.UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}

		//db operations
		err := toRun(client)

		if err != nil {
			err = sessionContext.AbortTransaction(sessionContext)
			return err
		}
		if err = sessionContext.CommitTransaction(sessionContext); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
