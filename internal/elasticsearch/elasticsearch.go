package elasticsearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/watchlist-kata/logger/internal/models"
	"strings"
)

// Client represents an Elasticsearch client.
type Client struct {
	client *elasticsearch.Client
}

// NewClient creates a new Elasticsearch client.
func NewClient(url string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Elasticsearch client: %v", err)
	}

	return &Client{client: client}, nil
}

// SendLog sends a log to Elasticsearch.
func (c *Client) SendLog(log models.Log) error {
	// Prepare the document for indexing
	doc := map[string]interface{}{
		"message":    log.Message,
		"@timestamp": log.Timestamp,
		"topic":      log.Topic,
	}
	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal log to JSON: %v", err)
	}

	// Send the document to Elasticsearch
	res, err := c.client.Index(
		"logs",
		strings.NewReader(string(docJSON)),
		c.client.Index.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("failed to send log to Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	// Check for errors in the response
	if res.IsError() {
		return fmt.Errorf("Elasticsearch error: %s", res.String())
	}

	return nil
}
