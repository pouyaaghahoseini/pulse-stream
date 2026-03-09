package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	EventsConsumedTotal      prometheus.Counter
	UnmarshalFailureTotal    prometheus.Counter
	DatabaseWriteFailureTotal prometheus.Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{
		EventsConsumedTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "analytics_events_consumed_total",
			Help: "Total number of events consumed by the analytics service.",
		}),
		UnmarshalFailureTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "analytics_unmarshal_failure_total",
			Help: "Total number of message unmarshal failures.",
		}),
		DatabaseWriteFailureTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "analytics_database_write_failure_total",
			Help: "Total number of failed database writes.",
		}),
	}

	prometheus.MustRegister(
		m.EventsConsumedTotal,
		m.UnmarshalFailureTotal,
		m.DatabaseWriteFailureTotal,
	)

	return m
}