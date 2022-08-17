package repository

import (
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserSet map[int]struct{}

var data = make(map[int64]UserSet)

type confirmRepo struct {
	client *mongo.Client
}

func (c confirmRepo) HasUsers(chatID int64) bool {
	_, ok := data[chatID]
	return ok
}

type ConfirmRepository interface {
	Add(chatID int64) error
	GetUsers(chatID int64) UserSet
	HasUsers(chatID int64) bool
	DeleteUser(chatID int64, userID int) error
	Finish(chatID int64)
}

func NewConfirmRepo(client *mongo.Client) ConfirmRepository {
	return confirmRepo{client: client}
}

func (c confirmRepo) Add(chatID int64) error {
	if _, ok := data[chatID]; ok {
		return nil
	}
	data[chatID] = make(UserSet)
	currentList := data[chatID]
	participantRepo := NewParticipantRepo(c.client)
	participants, err := participantRepo.GetByChatID(chatID)
	if err != nil {
		return err
	}
	for _, part := range participants {
		currentList[part.UserID] = struct{}{}
	}
	return nil
}

func (c confirmRepo) GetUsers(chatID int64) UserSet {
	if _, ok := data[chatID]; ok {
		return data[chatID]
	}
	return nil
}

func (c confirmRepo) DeleteUser(chatID int64, userID int) error {
	if _, ok := data[chatID]; !ok {
		return errors.New("chat not queued for confirmation")
	}
	if _, ok := data[chatID][userID]; ok {
		delete(data[chatID], userID)
	}
	return nil
}

func (c confirmRepo) Finish(chatID int64) {
	if _, ok := data[chatID]; !ok {
		return
	}
	delete(data, chatID)
}
