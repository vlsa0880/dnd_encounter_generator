package gin_middleware_prometheus

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type Config struct {
	Prometheus struct {
		Metrics struct {
			Endpoint string
		}
	}
}

func SetupPrometheusMetrics(router *gin.Engine, settingsLoader settings.SettingsLoader) error {
	config := Config{}
	if err := settingsLoader.Load(&config); err != nil {
		return err
	}

	router.Use(NewHttpRequestDurationMetric())
	router.Use(NewTotalHttpRequestMetric())

	router.GET(config.Prometheus.Metrics.Endpoint, gin.WrapH(promhttp.Handler()))
	return nil
}
