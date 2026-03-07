package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"pulse-stream/api-service/internal/kafka"
)

type PostEvent struct {
	PostID	string    `json:"post_id"`
	Platform string	`json:"platform"`
	ContentType string `json:"content_type"`
	EngagementScore int 	 `json:"engagement_score"`
	CreatedAt time.Time `json:"created_at"`
}

type App struct {
	producer *kafka.Producer
}

func (a *App) createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event PostEvent

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if event.PostID == "" || event.Platform == "" || event.ContentType == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if event.EngagementScore < 0 {
		http.Error(w, "Engagement score must be non-negative", http.StatusBadRequest)
		return
	}

	kafkaEvent := kafka.PostEvent{
		PostID: event.PostID,
		Platform: event.Platform,
		ContentType: event.ContentType,
		EngagementScore: event.EngagementScore,
		CreatedAt: event.CreatedAt,
	}

	err = a.producer.PublishPostEvent(context.Background(), kafkaEvent)
	if err != nil {
		log.Printf("failed to publish event: %v\n", err)
		http.Error(w, "failed to publish event", http.StatusInternalServerError)
		return
	}

	log.Printf("Received post event: %+v\n", event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "accepted and published"}`))

}

func main() {
	producer := kafka.NewProducer("localhost:9092", "post-events")
	defer producer.Close()
	
	app := &App{
		producer: producer,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/posts", app.createPostHandler)

	log.Println("API service is running on http://localhost:8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}