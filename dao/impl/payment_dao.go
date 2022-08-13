package impl

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"money-diff/dao/models"
	"time"
)

type PaymentDaoImpl struct {
	client *mongo.Client
}

func NewPaymentDao(client *mongo.Client) *PaymentDaoImpl {
	return &PaymentDaoImpl{client: client}
}

func (dao PaymentDaoImpl) Create(p *models.Payment) error {
	collection := dao.client.Database("money").Collection("payments")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := collection.InsertOne(ctx, p)
	fmt.Println(req.InsertedID)
	if err != nil {
		return fmt.Errorf("error inserting: %s", err)
	}

	return nil
}

func (dao PaymentDaoImpl) GetGroupedByChatID(chatID int64) ([]bson.M, error) {
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

	for _, result := range results {
		fmt.Println(result)
	}
	return results, nil
}

func (dao PaymentDaoImpl) GetByChatID(chatID int64) ([]models.Payment, error) {
	collection := dao.client.Database("money").Collection("payments")
	filter := bson.D{{"chat_id", chatID}}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}
	var results []models.Payment
	if err = cur.All(context.TODO(), &results); err != nil {
		return nil, fmt.Errorf("error querying: %s", err)
	}

	for _, result := range results {
		fmt.Println(result)
	}
	return results, nil
}
