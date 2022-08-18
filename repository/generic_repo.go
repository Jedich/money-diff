package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type genericRepo struct {
	client         *mongo.Client
	collectionName string
}

func (repo genericRepo) DeleteByChatID(ctx context.Context, chatID int64) error {
	collection := repo.client.Database("money").Collection(repo.collectionName)
	filter := bson.D{{"chat_id", chatID}}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting: %s", err)
	}
	return nil
}
