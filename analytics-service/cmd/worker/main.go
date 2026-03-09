package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"pulse-stream/analytics-service/internal/analytics"
	"pulse-stream/analytics-service/internal/store"
	"net/http"
	"pulse-stream/analytics-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	kafkago "github.com/segmentio/kafka-go"
)

type PostEvent struct {
	PostID          string    `json:"post_id"`
	Platform        string    `json:"platform"`
	ContentType     string    `json:"content_type"`
	EngagementScore int       `json:"engagement_score"`
	CreatedAt       time.Time `json:"created_at"`
}

func main() {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "post-events",
		GroupID: "analytics-service-group",
	})
	defer reader.Close()
	
	connectionString := "postgres://postgres:postgres@postgres:5432/pulsestream?sslmode=disable"
	dbStore, err := store.NewPostgresStore(connectionString)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer dbStore.Close()
	
	processor := analytics.NewProcessor()
	workerMetrics := metrics.NewMetrics()

	log.Println("Analytics service is listening for events on topic: post-events")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server running on http://localhost:8081/metrics")
		err := http.ListenAndServe(":8081", nil)
		if err != nil {
			log.Fatalf("failed to start metrics server: %v", err)
		}
	}()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("failed to read message: %v\n", err)
			continue
		}

		var event PostEvent
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			workerMetrics.UnmarshalFailureTotal.Inc()
			log.Printf("failed to unmarshal message: %v\n", err)
			continue
		}

		workerMetrics.EventsConsumedTotal.Inc()

		updatedStats := processor.ProcessEvent(event.Platform, event.EngagementScore)

		err = dbStore.UpsertPlatformStats(store.PlatformStats{
			Platform:        updatedStats.Platform,
			TotalPosts:      updatedStats.TotalPosts,
			TotalEngagement: updatedStats.TotalEngagement,
			AverageScore:    updatedStats.AverageScore,
		})
		if err != nil {
			workerMetrics.DatabaseWriteFailureTotal.Inc()
			log.Printf("failed to save stats to Postgres: %v\n", err)
			continue
		}

		log.Printf("processed and saved event: %+v\n", event)
		log.Printf(
			"saved stats: platform=%s total_posts=%d total_engagement=%d average_score=%.2f\n",
			updatedStats.Platform,
			updatedStats.TotalPosts,
			updatedStats.TotalEngagement,
			updatedStats.AverageScore,
		)
	}
}