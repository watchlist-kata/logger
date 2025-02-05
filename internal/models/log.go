package models

import "time"

type Log struct {
	Message   string    `json:"message"`
	Topic     string    `json:"topic"`
	Timestamp time.Time `json:"@timestamp"`
}
