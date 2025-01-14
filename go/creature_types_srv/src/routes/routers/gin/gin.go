package gin

import (
	"fmt"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	handlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin/handlers"
	gin_middleware "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin/middleware"
	gin_middleware_prometheus "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin/middleware/prometheus"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Config struct {
	Http struct {
		Address                  string
		Port                     string
		GetCreatureTypesEndpoint string
		PutCreatureTypesEndpoint string
	}
}

type GinManager struct {
	router      *gin.Engine
	dataHandler data_handlers_interfaces.DataHandler
	config      Config
}

func New(settingsLoader settings.SettingsLoader, dataHandler data_handlers_interfaces.DataHandler) *GinManager {
	mgr := GinManager{}
	if err := settingsLoader.Load(&mgr.config); err != nil {
		logger.GetInstance().Error(
			"Error loading gin router config",
			zap.Error(err),
		)
		return nil
	}
	if dataHandler == nil {
		logger.GetInstance().Error(
			"bad data handler",
		)
		return nil
	}
	mgr.dataHandler = dataHandler

	mgr.router = gin.Default()
	if err := mgr.setupMiddleware(settingsLoader); err != nil {
		logger.GetInstance().Error(
			"can't setup middleware",
			zap.Error(err),
		)
		return nil
	}

	mgr.setupRoutes(dataHandler)
	return &mgr
}

func (mgr *GinManager) Run() {
	if mgr.router == nil || mgr.dataHandler == nil {
		panic("Run Init() before Run()")
	}
	full_address := mgr.config.Http.Address + ":" + mgr.config.Http.Port
	if err := mgr.router.Run(full_address); err != nil {
		err_msg := fmt.Errorf("can't run gin: %s", err)
		panic(err_msg)
	}
}

func (mgr *GinManager) Stop() {}

func (mgr *GinManager) setupMiddleware(settingsLoader settings.SettingsLoader) error {
	mgr.router.Use(gin_middleware.NewLoggerHandler(settingsLoader))

	if err := gin_middleware_prometheus.SetupPrometheusMetrics(mgr.router, settingsLoader); err != nil {
		return err
	}

	return nil
}

func (mgr *GinManager) setupRoutes(dataHandler data_handlers_interfaces.DataHandler) {
	mgr.router.GET(
		mgr.config.Http.GetCreatureTypesEndpoint,
		handlers.GetCreatureTypesHandler(dataHandler),
	)
	mgr.router.PUT(
		mgr.config.Http.PutCreatureTypesEndpoint,
		handlers.PutCreatureTypesHandler(dataHandler),
	)
}
