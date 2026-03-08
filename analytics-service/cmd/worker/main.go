package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"pulse-stream/analytics-service/internal/analytics"
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

	processor := analytics.NewProcessor()

	log.Println("Analytics service is listening for events on topic: post-events")

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %v\n", err)
			continue
		}

		var event PostEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Printf("Failed to unmarshal message: %v\n", err)
			continue
		}

		processor.ProcessEvent(event.Platform, event.EngagementScore)

		log.Printf("processed event: %+v\n", event)
		log.Printf("total events processed: %d\n", processor.TotalEvents)

		for platform, stats := range processor.PlatformStats {
			log.Printf("Platform: %s, Total Posts: %d, Total Engagement: %d, Average Engagement: %.2f\n",
				platform, 
				stats.TotalPosts,
				stats.TotalEngagement, 
				stats.AverageEngagement,
			)
		}
	}
}