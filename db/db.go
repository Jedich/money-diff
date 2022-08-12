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

//OpenConnection return a connection of a desired db driver
func OpenConnection(ctx context.Context) *mongo.Client {
	//var dsn string
	//switch helpers.Getenv("DB_DRIVER") {
	//case "postgres":
	//	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=money port=%s sslmode=disable",
	//		helpers.Getenv("DB_HOST"), helpers.Getenv("DB_USER"), helpers.Getenv("DB_PASS"),
	//		helpers.Getenv("DB_PORT"))
	//default:
	//	log.Fatal("no driver configuration found")
	//}
	//
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	log.Fatal("failed to connect to database")
	//}
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
