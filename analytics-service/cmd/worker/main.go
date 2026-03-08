package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type PostEvent struct {
	PostID		  string    `json:"post_id"`
	Platform	  string    `json:"platform"`
	ContentType	  string    `json:"content_type"`
	EngagementScore int       `json:"engagement_score"`
	CreatedAt	  time.Time `json:"created_at"`
}

func main() {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "post-events",
		GroupID: "analytics-service-group",
	})

	defer reader.Close()

	log.Println("Analytics worker started, waiting for messages...")

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Error reading message: %v\n", err)
			continue
		}

		var event PostEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		log.Printf(
			"consumed message topic=%s partition=%d offset=%d event=%+v\n",
			message.Topic,
			message.Partition,
			message.Offset,
			event,
		)
	}
}