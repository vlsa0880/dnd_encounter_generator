package gin

import (
	"fmt"
	"net/http"
	"time"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	dataHandler data_handlers_interfaces.DataHandler
	config      Config
}

func New(settings_loader settings.SettingsLoader) *GinManager {
	mgr := GinManager{}
	mgr.router = gin.Default()
	mgr.setupMiddleware()
	if err := settings_loader.Load(&mgr.config); err != nil {
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
		err_msg := fmt.Errorf("can't run gin: %s", err)
		panic(err_msg)
	}
}

func (mgr *GinManager) Stop() {
}

func (mgr *GinManager) Init(data_handler data_handlers_interfaces.DataHandler) error {
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
	gin_ctx.IndentedJSON(http.StatusOK, mgr.dataHandler.GetCreatureTypes(gin_ctx))
}

func (mgr *GinManager) setupMiddleware() {
	mgr.router.Use(setupLogger())
}

func setupLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := time.Now()
		traceID := uuid.New().String()
		ctx.Set("traceID", traceID)

		ctx.Next()

		logger.GetInstance().Debug(
			"request processed data",
			zap.String("traceID", traceID),
			zap.Duration("processing_latency", time.Since(t)),
			zap.Int("http_status", ctx.Writer.Status()),
		)
	}
}
