package gin_middleware_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Setup(router *gin.Engine) {
	router.Use(NewHttpRequestDurationMetric())
	router.Use(NewTotalHttpRequestMetric())

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
