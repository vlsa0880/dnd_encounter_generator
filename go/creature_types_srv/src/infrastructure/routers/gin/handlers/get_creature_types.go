package gin_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	icontrollers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
)

func GetCreatureTypesHandler(controller icontrollers.CreatureTypes) gin.HandlerFunc {
	if controller == nil {
		panic("Bad data handler")
	}
	return func(ctx *gin.Context) {
		creatureTypes, err := controller.Get(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, entities.CreatureTypes{})
		} else {
			ctx.IndentedJSON(http.StatusOK, creatureTypes)
		}
	}
}
