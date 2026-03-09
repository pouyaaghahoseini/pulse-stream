package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	HTTPRequestsTotal       prometheus.Counter
	KafkaPublishSuccessTotal prometheus.Counter
	KafkaPublishFailureTotal prometheus.Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{
		HTTPRequestsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_http_requests_total",
			Help: "Total number of HTTP requests received by the API service.",
		}),
		KafkaPublishSuccessTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_kafka_publish_success_total",
			Help: "Total number of events successfully published to Kafka.",
		}),
		KafkaPublishFailureTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "api_kafka_publish_failure_total",
			Help: "Total number of failed Kafka publish attempts.",
		}),
	}

	prometheus.MustRegister(
		m.HTTPRequestsTotal,
		m.KafkaPublishSuccessTotal,
		m.KafkaPublishFailureTotal,
	)

	return m
}