package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoConn string

	KafkaConfig kafkaConfig
}

type kafkaConfig struct {
	Address      string
	KafkaTopic   string
	KafkaGroupID string
}

func NewConfig() (*Config, error) {
	// Load config.env from cwd or project root (when run via go run ./cmd/...)
	for _, name := range []string{"config.env", "./config.env", "../config.env"} {
		if err := godotenv.Load(name); err == nil {
			break
		}
	}
	// Optional: fail only if required vars are missing below

	config := Config{
		MongoConn: os.Getenv("DB_PATH"),
		KafkaConfig: kafkaConfig{
			Address:      os.Getenv("KafkaServerAddress"),
			KafkaTopic:   os.Getenv("KafkaTopic"),
			KafkaGroupID: os.Getenv("KafkaGroupID"),
		},
	}
	return &config, nil
}
