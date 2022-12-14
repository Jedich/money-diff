package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"money-diff/model"
	"time"
)

type ParticipantRepository interface {
	Create(ctx context.Context, p *model.Participant) error
	GetByChatID(chatID int64) ([]model.Participant, error)
	DeleteByChatID(ctx context.Context, chatID int64) error
}

type participantRepoImpl struct {
	*genericRepo
}

func NewParticipantRepo(client *mongo.Client) ParticipantRepository {
	return participantRepoImpl{&genericRepo{client: client, collectionName: "participants"}}
}

func (repo participantRepoImpl) Create(ctx context.Context, p *model.Participant) error {
	collection := repo.client.Database("money").Collection("participants")

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		if e, ok := err.(mongo.WriteException); ok {
			for _, we := range e.WriteErrors {
				if we.Code == 11000 {
					log.Println("Found duplicate in participants, skipping")
					return nil
				}
			}
		}
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (repo participantRepoImpl) GetByChatID(chatID int64) ([]model.Participant, error) {
	collection := repo.client.Database("money").Collection("participants")
	filter := bson.D{{"chat_id", chatID}}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}
	var results []model.Participant
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	return results, nil
}
