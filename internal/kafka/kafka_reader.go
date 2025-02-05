package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/watchlist-kata/logger/internal/elasticsearch"
	"github.com/watchlist-kata/logger/internal/models"
	"log/slog"
	"time"
)

// StartKafkaReader starts a Kafka reader for the specified topic.
func StartKafkaReader(topic string, brokers []string, elasticsearchURL string, logger *slog.Logger) {
	// Log the start of the Kafka reader
	logger.Info("Starting Kafka reader", "topic", topic)

	// Create a Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		Partition:      0,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: 0,
	})
	defer reader.Close()

	// Create an Elasticsearch client
	esClient, err := elasticsearch.NewClient(elasticsearchURL)
	if err != nil {
		logger.Error("Failed to create Elasticsearch client", "error", err)
		return
	}

	// Infinite loop to read messages
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.Error("Failed to read from Kafka", "topic", topic, "error", err)
			continue
		}

		logger.Info("Received log", "topic", topic, "message", string(msg.Value))

		// Create a log entry
		logEntry := models.Log{
			Message:   string(msg.Value),
			Topic:     topic,
			Timestamp: time.Now(),
		}

		// Send the log to Elasticsearch
		err = esClient.SendLog(logEntry)
		if err != nil {
			logger.Error("Failed to send log to Elasticsearch", "topic", topic, "error", err)
		} else {
			logger.Info("Log successfully sent to Elasticsearch", "topic", topic)
		}
	}
}
