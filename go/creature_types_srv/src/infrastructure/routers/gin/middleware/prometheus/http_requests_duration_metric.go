package gin_middleware_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func NewHttpRequestDurationMetric() gin.HandlerFunc {
	httpRequestDuration := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
	prometheus.MustRegister(httpRequestDuration)
	return func(ctx *gin.Context) {
		timer := prometheus.NewTimer(httpRequestDuration)
		defer timer.ObserveDuration()
		ctx.Next()
	}
}
