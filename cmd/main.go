package main

import (
	"context"
	consumer "gw-notification/internal/broker/kafka"
	"gw-notification/internal/config"
	"gw-notification/internal/logger"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	logger, err := logger.InitLogger("gw-notification")
	if err != nil {
		log.Fatal("Failed to initialize logger", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("Error reading config", "error", err)
		return
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConn))
	if err != nil {
		logger.Error("Error connection to mongoDB", "error", err)
	}
	db := client.Database("notifier_db")

	go func() {
		sig := <-signals
		logger.Info("Received termination signal", "signal", sig.String())
		cancel()
	}()

	var wg sync.WaitGroup
	numWorkers := 5

	logger.Info("Starting worker pool", "workers_count", numWorkers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			consumer.StartConsumer(i, ctx, cfg, db, logger)
		}(i)
	}

	wg.Wait()

}
