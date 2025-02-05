package main

import (
	"github.com/watchlist-kata/logger/internal/config"
	"github.com/watchlist-kata/logger/internal/kafka"
	"log/slog"
	"os"
	"sync"
)

func main() {
	// Set up the logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	logger.Info("Logger started", "topics", cfg.KafkaTopics)

	// Create a WaitGroup to synchronize goroutine
	var wg sync.WaitGroup

	// Start Kafka readers for each topic
	for _, topic := range cfg.KafkaTopics {
		wg.Add(1)
		go func(topic string) {
			defer wg.Done()
			kafka.StartKafkaReader(topic, cfg.KafkaBrokers, cfg.ElasticsearchURL, logger)

		}(topic)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
