package gin_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type GetHandler struct {
	creatureTypesDB dbinterfaces.CreatureTypes
}

func GetCreatureTypesHandler(creatureTypesDB dbinterfaces.CreatureTypes) gin.HandlerFunc {
	if creatureTypesDB == nil {
		panic("Bad data handler")
	}
	handler := GetHandler{}
	handler.creatureTypesDB = creatureTypesDB
	return func(ctx *gin.Context) {
		creatureTypes, err := handler.creatureTypesDB.GetCreatureTypes(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, entities.CreatureTypes{})
		} else {
			ctx.IndentedJSON(http.StatusOK, creatureTypes)
		}
	}
}
