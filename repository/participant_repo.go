package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/model"
	"time"
)

type ParticipantRepository interface {
	Create(p *model.Participant) error
	GetByChatID(chatID int64) ([]model.Participant, error)
}

type ParticipantRepoImpl struct {
	client *mongo.Client
}

func NewParticipantRepo(client *mongo.Client) ParticipantRepository {
	return ParticipantRepoImpl{client: client}
}

func (dao ParticipantRepoImpl) Create(p *model.Participant) error {
	collection := dao.client.Database("money").Collection("participants")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (dao ParticipantRepoImpl) GetByChatID(chatID int64) ([]model.Participant, error) {
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
