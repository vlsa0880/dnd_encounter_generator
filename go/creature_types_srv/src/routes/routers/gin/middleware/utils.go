package gin_middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"go.uber.org/zap"
)

type Request struct {
	TraceIDField string
}

func getTraceIDField(settingsLoader settings.SettingsLoader) string {
	requestConfig := Request{}
	if err := settingsLoader.Load(&requestConfig); err != nil {
		return "traceID"
	}
	return requestConfig.TraceIDField
}

func getTraceID(ctx *gin.Context, traceIDField string) string {
	traceID, exist := ctx.Get(traceIDField)
	if !exist {
		traceID = uuid.New().String()
		ctx.Set(traceIDField, traceID)
	}
	resValue, ok := traceID.(string)
	if !ok {
		resValue = uuid.New().String()
		ctx.Set(traceIDField, resValue)

		logger.GetInstance().Error(
			"traceID must be string. Updating it",
			zap.String("traceid_value", fmt.Sprintf("%v", traceID)),
			zap.String("traceid_type", fmt.Sprintf("%T", traceID)),
			zap.String("updated_traceid", resValue),
		)
	}
	return resValue
}
