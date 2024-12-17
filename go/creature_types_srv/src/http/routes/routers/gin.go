package routers

import (
	"creature_types_srv/src/data_handlers/idata_handler"
	logger "creature_types_srv/src/logger/zap"
	settings "creature_types_srv/src/settings/loader"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	Http struct {
		Address            string
		Port               string
		GetCreatureTypesEP string
	}
}

type GinManager struct {
	router      *gin.Engine
	dataHandler idata_handler.IDataHandler
	config      *Config
}

func New(settings_loader settings.ISettingsLoader) *GinManager {
	mgr := GinManager{}
	mgr.router = gin.Default()
	mgr.setupMiddleware()
	mgr.config = &Config{}
	if err := settings_loader.Load(mgr.config); err != nil {
		logger.GetInstance().Error(
			"Error loading gin router config",
			zap.String("msg", err.Error()),
		)
		return nil
	}
	return &mgr
}

func (mgr *GinManager) Run() {
	if mgr.router == nil || mgr.dataHandler == nil {
		panic("Run Init() and SetupDataHandler() before Run()")
	}
	full_address := mgr.config.Http.Address + ":" + mgr.config.Http.Port
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
	mgr.dataHandler = data_handler
	mgr.router.GET(
		mgr.config.Http.GetCreatureTypesEP,
		mgr.setupGetCreatureTypes,
	)
	return nil
}

func (mgr *GinManager) setupGetCreatureTypes(gin_ctx *gin.Context) {
	gin_ctx.IndentedJSON(http.StatusOK, mgr.dataHandler.GetCreatureTypes())
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
