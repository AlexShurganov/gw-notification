package mongo

import (
	"context"
	"fmt"
	"gw-notification/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	collection *mongo.Collection
}

func NewStorage(db *mongo.Database) *MongoStorage {

	return &MongoStorage{
		collection: db.Collection("wallet_operations"),
	}

}
func (s *MongoStorage) StoreOperation(walletID, amount string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var notification models.Notification
	notification.WalletID = walletID
	notification.Amount = amount
	notification.Timestamp = time.Now()

	_, err := s.collection.InsertOne(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to insert operation into MongoDB: %w", err)
	}

	return nil
}
