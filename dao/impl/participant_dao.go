package impl

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"money-diff/dao/models"
	"time"
)

type ParticipantDaoImpl struct {
	client *mongo.Client
}

func NewParticipantDao(client *mongo.Client) *ParticipantDaoImpl {
	return &ParticipantDaoImpl{client: client}
}

func (dao ParticipantDaoImpl) Create(p *models.Participant) error {
	collection := dao.client.Database("money").Collection("participants")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", p.ID}}
	update := bson.D{{"$set", bson.D{{"chat_id", p.ChatID}}}, {"$set", bson.D{{"user_id", p.UserID}}}}

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (dao ParticipantDaoImpl) GetByChatID(chatID int64) ([]models.Participant, error) {
	collection := dao.client.Database("money").Collection("participants")
	filter := bson.D{{"chat_id", chatID}}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}
	var results []models.Participant
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	return results, nil
}
