package gin_middleware

import (
	"time"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func NewLoggerHandler(settings_loader settings.SettingsLoader) gin.HandlerFunc {
	traceIDField := getTraceIDField(settings_loader)
	return func(ctx *gin.Context) {
		t := time.Now()
		logger.GetInstance().Info(
			"request processing started",
			zap.String("traceID", getTraceID(ctx, traceIDField)),
			zap.Duration("processing_latency", time.Since(t)),
			zap.Int("http_status", ctx.Writer.Status()),
		)

		ctx.Next()

		logger.GetInstance().Info(
			"request processing finished",
			zap.String("traceID", getTraceID(ctx, traceIDField)),
			zap.Duration("processing_latency", time.Since(t)),
			zap.Int("http_status", ctx.Writer.Status()),
		)
	}
}
