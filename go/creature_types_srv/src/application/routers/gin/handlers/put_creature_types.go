package gin_handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	json_validators "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/gin/validators/json"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type PutHandler struct {
	creatureTypesDB dbinterfaces.CreatureTypes
}

type jsonTypedData struct {
	Types []string
}

func PutCreatureTypesHandler(creatureTypesDB dbinterfaces.CreatureTypes) gin.HandlerFunc {
	if creatureTypesDB == nil {
		panic("Bad creature types db")
	}
	handler := PutHandler{}
	handler.creatureTypesDB = creatureTypesDB
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
		creatureTypesData := entities.CreatureTypes{}
		for _, typeName := range typesData.Types {
			creatureTypesData.Types = append(creatureTypesData.Types, entities.CreatureType{Name: typeName})
		}
		if err := handler.creatureTypesDB.UploadCreatureTypes(ctx, &creatureTypesData); err != nil {
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
