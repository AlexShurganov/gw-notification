package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"gw-notification/internal/config"
	mongodb "gw-notification/internal/storage/mongo"
	"gw-notification/pkg/models"
	"log/slog"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func StartConsumer(id int, ctx context.Context, cfg *config.Config, db *mongo.Database, logger *slog.Logger) error {
	l := logger.With("worker_id", id)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.KafkaConfig.Address},
		Topic:   cfg.KafkaConfig.KafkaTopic,
		GroupID: cfg.KafkaConfig.KafkaGroupID,
	})
	defer r.Close()

	st := mongodb.NewStorage(db)
	l.Info("Consumer worker ready", "topic", cfg.KafkaConfig.KafkaTopic)
	for {
		select {
		case <-ctx.Done():
			l.Info("worker stopping", "info", id)
			return nil
		default:
			msg, err := r.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || ctx.Err() != nil {
					l.Debug("Worker context canceled, exiting loop")
					return nil
				}
				l.Error("Kafka read error", "error", err)
				return err
			}

			var notification models.Notification
			if err := json.Unmarshal(msg.Value, &notification); err != nil {
				l.Error("Error unmarshaling notification", "error", err)
				continue
			}

			if err = st.StoreOperation(notification.WalletID, notification.Amount); err != nil {
				l.Error("Error storing operation", "error", err)
			}

			l.Info("Notification saved", "offset", msg.Offset, "wallet_id", notification.WalletID)
		}

	}
}
