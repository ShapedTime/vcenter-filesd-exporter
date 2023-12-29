package promhelper

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	errors               *prometheus.CounterVec
	promRequestsReceived *prometheus.CounterVec
)

func init() {
	errors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "errors",
		Help: "Number of errors",
	}, []string{})

	promRequestsReceived = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "prom_requests_received",
		Help: "Number of requests received on prom endpoint",
	}, []string{"path"})
}

func Error() {
	errors.WithLabelValues().Inc()
}

func RequestReceived(path string) {
	promRequestsReceived.WithLabelValues(path).Inc()
}
