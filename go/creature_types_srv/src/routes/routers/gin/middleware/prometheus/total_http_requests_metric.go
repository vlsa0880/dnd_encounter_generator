package gin_middleware_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func NewTotalHttpRequestMetric() gin.HandlerFunc {
	httpRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	prometheus.MustRegister(httpRequests)
	return func(ctx *gin.Context) {
		ctx.Next()
		httpRequests.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Inc()
	}
}
