package models

import "time"

type Notification struct {
	WalletID  string    `json:"wallet_id" bson:"wallet_id"`
	Amount    string    `json:"amount" bson:"amount"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
