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
}

type participantRepoImpl struct {
	client *mongo.Client
}

func NewParticipantRepo(client *mongo.Client) ParticipantRepository {
	return participantRepoImpl{client: client}
}

func (dao participantRepoImpl) Create(ctx context.Context, p *model.Participant) error {
	collection := dao.client.Database("money").Collection("participants")

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

func (dao participantRepoImpl) GetByChatID(chatID int64) ([]model.Participant, error) {
	collection := dao.client.Database("money").Collection("participants")
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
