package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type PostEvent struct {
	PostID	string    `json:"post_id"`
	Platform string	`json:"platform"`
	ContentType string `json:"content_type"`
	EngagementScore int 	 `json:"engagement_score"`
	CreatedAt time.Time `json:"created_at"`
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	log.Printf("Received post event: %+v\n", event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "accepted"}`))

}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/posts", createPostHandler)

	log.Println("API service is running on http://localhost:8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}