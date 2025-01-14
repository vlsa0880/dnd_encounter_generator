package gin_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	creature_types_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
)

type GetHandler struct {
	dataHandler data_handlers_interfaces.DataHandler
}

func GetCreatureTypesHandler(dataHandler data_handlers_interfaces.DataHandler) gin.HandlerFunc {
	if dataHandler == nil {
		panic("Bad data handler")
	}
	handler := GetHandler{}
	handler.dataHandler = dataHandler
	return func(ctx *gin.Context) {
		creatureTypes, err := handler.dataHandler.GetCreatureTypes(ctx)
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, creature_types_data.Data{})
		} else {
			ctx.IndentedJSON(http.StatusOK, creatureTypes)
		}
	}
}
