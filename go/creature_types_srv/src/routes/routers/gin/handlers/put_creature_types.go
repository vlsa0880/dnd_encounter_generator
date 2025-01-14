package gin_handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	creature_types_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	json_validators "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin/validators/json"
)

type PutHandler struct {
	dataHandler data_handlers_interfaces.DataHandler
}

type jsonTypedData struct {
	Types []string
}

func PutCreatureTypesHandler(dataHandler data_handlers_interfaces.DataHandler) gin.HandlerFunc {
	if dataHandler == nil {
		panic("Bad data handler")
	}
	handler := PutHandler{}
	handler.dataHandler = dataHandler
	return func(ctx *gin.Context) {
		if err := json_validators.JsonBodyExist(ctx); err != nil {
			logger.GetInstance().Error(
				"bad json body",
				zap.Error(err),
			)
			ctx.Status(http.StatusBadRequest)
			return
		}
		typesData, err := jsonMapping(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "bad json data")
			return
		}
		creatureTypesData := creature_types_data.Data{}
		for _, typeName := range typesData.Types {
			creatureTypesData.Types = append(creatureTypesData.Types, creature_types_data.Type{Name: typeName})
		}
		if err := handler.dataHandler.PutCreatureTypes(ctx, &creatureTypesData); err != nil {
			logger.GetInstance().Error(
				"can't handle request",
				zap.Error(err),
			)
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusOK)
	}
}

func jsonMapping(ctx *gin.Context) (*jsonTypedData, error) {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.GetInstance().Error(
			"can't read json body",
			zap.Error(err),
		)
		return nil, err
	}
	typesData := jsonTypedData{}
	if err = json.Unmarshal(jsonData, &typesData); err != nil {
		logger.GetInstance().Error(
			"can't unmarshal json",
			zap.String("json", string(jsonData)),
			zap.Error(err),
		)
		return nil, err
	}
	return &typesData, nil
}
