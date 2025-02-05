package gin_validators_json

import (
	"net/http"

	"github.com/gin-gonic/gin"
	jsonerrors "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/gin/errors/json"
)

func JsonBodyExist(ctx *gin.Context) error {
	if ctx.Request.Body == http.NoBody {
		return &jsonerrors.NoJsonBody{}
	}
	return nil
}
