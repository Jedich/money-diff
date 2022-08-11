package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"money-diff/bot/helpers"
)

//OpenConnection return a connection of a desired db driver
func OpenConnection() *mongo.Client {
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
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?maxPoolSize=20&w=majority",
		helpers.Getenv("DB_HOST"), helpers.Getenv("DB_USER"), helpers.Getenv("DB_PASS"))

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	return client
}
