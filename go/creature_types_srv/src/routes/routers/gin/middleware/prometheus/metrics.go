package gin_middleware_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Setup(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Use(NewTotalHttpRequestMetric())
}
