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
	GetGroupByChatID(chatID int64) ([]model.GroupedPayment, error)
	GetByChatID(chatID int64) ([]model.Payment, error)
	DeleteByChatID(ctx context.Context, chatID int64) error
}

type paymentRepoImpl struct {
	*genericRepo
}

func NewPaymentRepo(client *mongo.Client) PaymentRepository {
	return paymentRepoImpl{&genericRepo{client: client, collectionName: "payments"}}
}

func (repo paymentRepoImpl) Create(ctx context.Context, p *model.Payment) error {
	collection := repo.client.Database("money").Collection("payments")

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (repo paymentRepoImpl) GetGroupByChatID(chatID int64) ([]model.GroupedPayment, error) {
	collection := repo.client.Database("money").Collection("payments")
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
	var results []model.GroupedPayment
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	return results, nil
}

func (repo paymentRepoImpl) GetByChatID(chatID int64) ([]model.Payment, error) {
	collection := repo.client.Database("money").Collection("payments")
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
