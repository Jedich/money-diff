package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/model"
	"time"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *model.Payment) error
	GetGroupByChatID(chatID int64) ([]bson.M, error)
	GetByChatID(chatID int64) ([]model.Payment, error)
}

type paymentRepoImpl struct {
	client *mongo.Client
}

func NewPaymentRepo(client *mongo.Client) PaymentRepository {
	return paymentRepoImpl{client: client}
}

func (dao paymentRepoImpl) Create(ctx context.Context, p *model.Payment) error {
	collection := dao.client.Database("money").Collection("payments")

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (dao paymentRepoImpl) GetGroupByChatID(chatID int64) ([]bson.M, error) {
	collection := dao.client.Database("money").Collection("payments")
	filter := bson.D{
		{"$match", bson.D{{"chat_id", chatID}}}}
	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$username"},
			{"value", bson.D{{"$sum", "$value"}}},
		}}}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := collection.Aggregate(ctx, mongo.Pipeline{filter, groupStage})
	if err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}
	var results []bson.M
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	return results, nil
}

func (dao paymentRepoImpl) GetByChatID(chatID int64) ([]model.Payment, error) {
	collection := dao.client.Database("money").Collection("payments")
	filter := bson.D{{"chat_id", chatID}}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}
	var results []model.Payment
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	return results, nil
}
