package config

import (
	"errors"
	"os"
)

type Config struct {
	KafkaBrokers     []string
	KafkaTopics      []string
	ElasticsearchURL string
}

func LoadConfig() (*Config, error) {
	// Load configuration from environment variables
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return nil, errors.New("KAFKA_BROKERS not set")
	}

	elasticsearchURL := os.Getenv("ELASTICSEARCH_URL")
	if elasticsearchURL == "" {
		return nil, errors.New("ELASTICSEARCH_URL not set")
	}

	// List of Kafka topics to consume
	kafkaTopics := []string{"auth", "media", "subscription", "user", "gateway", "reviews", "watchlist"}

	return &Config{
		KafkaBrokers:     []string{kafkaBrokers},
		KafkaTopics:      kafkaTopics,
		ElasticsearchURL: elasticsearchURL,
	}, nil
}
