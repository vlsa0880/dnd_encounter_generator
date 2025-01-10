package gin

import (
	"fmt"
	"net/http"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	gin_middleware "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin/middleware"
	gin_middleware_prometheus "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin/middleware/prometheus"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

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
	dataHandler data_handlers_interfaces.DataHandler
	config      Config
}

func New(settingsLoader settings.SettingsLoader, dataHandler data_handlers_interfaces.DataHandler) *GinManager {
	mgr := GinManager{}
	if err := settingsLoader.Load(&mgr.config); err != nil {
		logger.GetInstance().Error(
			"Error loading gin router config",
			zap.String("msg", err.Error()),
		)
		return nil
	}
	if err := mgr.setupDataHandler(dataHandler); err != nil {
		logger.GetInstance().Error(
			"can't setup data handler",
			zap.String("msg", err.Error()),
		)
		return nil
	}
	mgr.router = gin.Default()
	if err := mgr.setupMiddleware(settingsLoader); err != nil {
		logger.GetInstance().Error(
			"can't setup middleware",
			zap.String("msg", err.Error()),
		)
		return nil
	}
	mgr.router.GET(
		mgr.config.Http.GetCreatureTypesEP,
		mgr.setupGetCreatureTypes,
	)
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

func (mgr *GinManager) setupDataHandler(dataHandler data_handlers_interfaces.DataHandler) error {
	if dataHandler == nil {
		return fmt.Errorf("data handler is nil")
	}
	mgr.dataHandler = dataHandler
	return nil
}

func (mgr *GinManager) setupGetCreatureTypes(ginCtx *gin.Context) {
	ginCtx.IndentedJSON(http.StatusOK, mgr.dataHandler.GetCreatureTypes(ginCtx))
}

func (mgr *GinManager) setupMiddleware(settingsLoader settings.SettingsLoader) error {
	mgr.router.Use(gin_middleware.NewLoggerHandler(settingsLoader))

	if err := gin_middleware_prometheus.SetupPrometheusMetrics(mgr.router, settingsLoader); err != nil {
		return err
	}
	return nil
}
