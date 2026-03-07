## Project Summary 
Pulse Stream is a backend system that receives social media post events, publishes them to Kafka, and processes them asynchronously to generate analytics.

## Main Components 
- API Service: accepts post events over HTTP
- Kafka: transports events between services
- Analytics Service: consumes events and calculates analytics

## Event Schema
An example of an even json:
{
  "post_id": "p_001",
  "platform": "twitter",
  "content_type": "text",
  "engagement_score": 42,
  "created_at": "2026-03-06T10:00:00Z"
}

## Processing Flow
	1.	Client sends a post event to the API
	2.	API validates the data
	3.	API publishes the event to Kafka
	4.	Analytics service reads the event from Kafka
	5.	Analytics service updates simple analytics