package routers

import (
	"creature_types_srv/src/data_handlers/idata_handler"
	logger "creature_types_srv/src/logger/zap"
	settings "creature_types_srv/src/settings/implementations"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GinManager struct {
	router       *gin.Engine
	data_handler idata_handler.IDataHandler
}

func (mgr *GinManager) Init() error {
	mgr.router = gin.Default()
	mgr.setupMiddleware()
	return nil
}

func (mgr *GinManager) Run() {
	if mgr.router == nil || mgr.data_handler == nil {
		panic("Run Init() and SetupDataHandler() before Run()")
	}
	full_address := settings.GetInstance().GetSettings().Http.Address + ":" + settings.GetInstance().GetSettings().Http.Port
	if err := mgr.router.Run(full_address); err != nil {
		panic(err)
	}
}

func (mgr *GinManager) SetupDataHandler(data_handler idata_handler.IDataHandler) error {
	if mgr.router == nil {
		return fmt.Errorf("gin not inited - gin.Engine is nil")
	}
	if data_handler == nil {
		return fmt.Errorf("data handler is nil")
	}
	mgr.data_handler = data_handler
	mgr.router.GET(
		settings.GetInstance().GetSettings().CreatureTypes.GetEndpoint,
		mgr.setupGetCreatureTypes,
	)
	return nil
}

func (mgr *GinManager) setupGetCreatureTypes(gin_ctx *gin.Context) {
	gin_ctx.IndentedJSON(http.StatusOK, mgr.data_handler.GetCreatureTypes())
}

func (mgr *GinManager) setupMiddleware() {
	mgr.router.Use(setupLogger())
}

func setupLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()

		ctx.Next()

		logger.GetInstance().Debug(
			"request processing data",
			zap.Duration("processing_latency", time.Since(t)),
			zap.Int("http_status", ctx.Writer.Status()),
		)
	}
}
