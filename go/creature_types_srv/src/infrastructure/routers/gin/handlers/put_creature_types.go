package gin_handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	icontrollers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	jsonvalidators "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/gin/validators/json"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
)

type jsonTypedData struct {
	Types []string
}

func PutCreatureTypesHandler(controller icontrollers.CreatureTypes) gin.HandlerFunc {
	if controller == nil {
		panic("Bad controller")
	}
	return func(ctx *gin.Context) {
		if err := jsonvalidators.JsonBodyExist(ctx); err != nil {
			logger.GetInstance().Error(
				"bad json body",
				zap.Error(err),
			)
			ctx.Status(http.StatusBadRequest)
			return
		}
		typesData := jsonTypedData{}
		err := jsonMapping(ctx, &typesData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "bad json data")
			return
		}
		creatureTypesData := entities.CreatureTypes{
			Types: make([]entities.CreatureType, 0, len(typesData.Types)),
		}
		for _, typeName := range typesData.Types {
			creatureTypesData.Types = append(creatureTypesData.Types, entities.CreatureType{Name: typeName})
		}
		if err := controller.Set(ctx, &creatureTypesData); err != nil {
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

func jsonMapping(ctx *gin.Context, typesData *jsonTypedData) error {
	jsonData, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.GetInstance().Error(
			"can't read json body",
			zap.Error(err),
		)
		return err
	}
	if err = json.Unmarshal(jsonData, typesData); err != nil {
		logger.GetInstance().Error(
			"can't unmarshal json",
			zap.String("json", string(jsonData)),
			zap.Error(err),
		)
		return err
	}
	return nil
}
