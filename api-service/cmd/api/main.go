package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"pulse-stream/api-service/internal/kafka"
	"pulse-stream/api-service/internal/store"
	"pulse-stream/api-service/internal/validation"
	"pulse-stream/api-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PostEvent struct {
	PostID          string    `json:"post_id"`
	Platform        string    `json:"platform"`
	ContentType     string    `json:"content_type"`
	EngagementScore int       `json:"engagement_score"`
	CreatedAt       time.Time `json:"created_at"`
}

type App struct {
	producer *kafka.Producer
	store    *store.PostgresStore
	metrics  *metrics.Metrics
}

func (a *App) createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	a.metrics.HTTPRequestsTotal.Inc()

	var event PostEvent

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err = validation.ValidatePostEvent(validation.PostEvent{
		PostID:          event.PostID,
		Platform:        event.Platform,
		ContentType:     event.ContentType,
		EngagementScore: event.EngagementScore,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	kafkaEvent := kafka.PostEvent{
		PostID:          event.PostID,
		Platform:        event.Platform,
		ContentType:     event.ContentType,
		EngagementScore: event.EngagementScore,
		CreatedAt:       event.CreatedAt,
	}

	err = a.producer.PublishPostEvent(context.Background(), kafkaEvent)
	if err != nil {
		a.metrics.KafkaPublishFailureTotal.Inc()
		log.Printf("failed to publish event: %v\n", err)
		http.Error(w, "failed to publish event", http.StatusInternalServerError)
		return
	}

	a.metrics.KafkaPublishSuccessTotal.Inc()
	log.Printf("published event to Kafka: %+v\n", event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted and published"}`))
}

func (a *App) getPlatformAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	results, err := a.store.GetAllPlatformAnalytics()
	if err != nil {
		log.Printf("failed to fetch analytics: %v\n", err)
		http.Error(w, "failed to fetch analytics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Printf("failed to encode analytics response: %v\n", err)
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	producer := kafka.NewProducer("localhost:9092", "post-events")
	defer producer.Close()

	connectionString := "postgres://postgres:postgres@localhost:5432/pulsestream?sslmode=disable"
	dbStore, err := store.NewPostgresStore(connectionString)
	if err != nil {
		log.Fatalf("failed to connect to Postgres: %v", err)
	}
	defer dbStore.Close()

	appMetrics := metrics.NewMetrics()

	app := &App{
		producer: producer,
		store:    dbStore,
		metrics:  appMetrics,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/posts", app.createPostHandler)
	mux.HandleFunc("/analytics/platforms", app.getPlatformAnalyticsHandler)
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("API service running on http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}