package instrument

import (
	"ecommerce-api/internal/config"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type RequestInstrument struct {
	RequestCount   *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
}

func NewRequestInstrument(cfg config.AppConfig) *RequestInstrument {
	labels := []string{"op"}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: fmt.Sprintf("%s_request_count", cfg.Name),
		Help: "Total request count",
	}, labels)

	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: fmt.Sprintf("%s_request_latency", cfg.Name),
		Help: "Total duration of request in microseconds",
	}, labels)

	prometheus.MustRegister(counter, latency)

	return &RequestInstrument{
		RequestCount:   counter,
		RequestLatency: latency,
	}
}
