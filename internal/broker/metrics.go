package broker

import "github.com/prometheus/client_golang/prometheus"

var(

	MessagesPublished = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "raven_messages_published_total",
		Help: "Total number of messages published",
	})

	MessagesConsumed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "raven_messages_consumed_total",
		Help: "Total number of messages consumed",
	})

	ReplicationErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "raven_replication_errors_total",
		Help: "Total number of replication errors",
	})
)
// Here it registers the metrics with Prometheus — tells Prometheus "these metrics exist, track them." 
// Without registration, Prometheus doesn't know about them and won't expose them on the /metrics endpoint.
func init(){
	prometheus.MustRegister(MessagesPublished)
	prometheus.MustRegister(MessagesConsumed)
	prometheus.MustRegister(ReplicationErrors)
}