package kafka

import (
	"context"
	"encoding/json"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type PostEvent struct {
	PostID          string    `json:"post_id"`
	Platform        string    `json:"platform"`
	ContentType     string    `json:"content_type"`
	EngagementScore int       `json:"engagement_score"`
	CreatedAt       time.Time `json:"created_at"`
}

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokerAddress string, topic string) *Producer {
	writer := &kafkago.Writer{
		Addr:     kafkago.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafkago.LeastBytes{},
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) PublishPostEvent(ctx context.Context, event PostEvent) error {
	messageBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	message := kafkago.Message{
		Key:  []byte(event.PostID),
		Value: messageBytes,
	}

	return p.writer.WriteMessages(ctx, message)
}

func (p *Producer) Close() error {
	return p.writer.Close()
}