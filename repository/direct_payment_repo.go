package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/model"
	"time"
)

type DirectPaymentRepository interface {
	Create(ctx context.Context, p *model.DirectPayment) error
	GetGroupByChatID(chatID int64) ([]bson.M, error)
	GetByChatID(chatID int64) ([]model.DirectPayment, error)
}

type directPaymentRepoImpl struct {
	client *mongo.Client
}

func NewDirectPaymentRepo(client *mongo.Client) DirectPaymentRepository {
	return directPaymentRepoImpl{client: client}
}

func (dao directPaymentRepoImpl) Create(ctx context.Context, p *model.DirectPayment) error {
	collection := dao.client.Database("money").Collection("direct_payments")

	_, err := collection.InsertOne(ctx, p)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (dao directPaymentRepoImpl) GetGroupByChatID(chatID int64) ([]bson.M, error) {
	collection := dao.client.Database("money").Collection("direct_payments")
	filter := bson.D{
		{"$match", bson.D{{"chat_id", chatID}}}}
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				bson.E{Key: "from", Value: "$from"},
				bson.E{Key: "to", Value: "$to"},
			}},
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

func (dao directPaymentRepoImpl) GetByChatID(chatID int64) ([]model.DirectPayment, error) {
	collection := dao.client.Database("money").Collection("direct_payments")
	filter := bson.D{{"chat_id", chatID}}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}
	var results []model.DirectPayment
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	return results, nil
}
