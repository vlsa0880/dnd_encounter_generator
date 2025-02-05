package ginrouter

import (
	"fmt"

	icontrollers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/interfaces"
	handlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/gin/handlers"
	gin_middleware "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/gin/middleware"
	gin_middleware_prometheus "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/gin/middleware/prometheus"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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
	router     *gin.Engine
	config     Config
	controller icontrollers.CreatureTypes
}

func New(settingsLoader settings.SettingsLoader, controller icontrollers.CreatureTypes) (*GinManager, error) {
	mgr := GinManager{}
	if err := settingsLoader.Load(&mgr.config); err != nil {
		return nil, fmt.Errorf("Error loading gin router config: %w", err)
	}
	if controller == nil {
		return nil, fmt.Errorf("bad data handler")
	}
	mgr.controller = controller

	mgr.router = gin.New()
	if err := mgr.setupMiddleware(settingsLoader); err != nil {
		return nil, fmt.Errorf("can't setup middleware: %w", err)
	}

	mgr.setupRoutes(controller)
	return &mgr, nil
}

func (mgr *GinManager) Run() error {
	if mgr.router == nil {
		return fmt.Errorf("Bad router")
	}
	if mgr.controller == nil {
		return fmt.Errorf("Bad db")
	}
	full_address := mgr.config.Http.Address + ":" + mgr.config.Http.Port
	if err := mgr.router.Run(full_address); err != nil {
		return fmt.Errorf("can't run gin: %s", err)
	}
	return nil
}

func (mgr *GinManager) Stop() error {
	return nil
}

func (mgr *GinManager) setupMiddleware(settingsLoader settings.SettingsLoader) error {
	mgr.router.Use(
		gin.Recovery(),
		gin_middleware.NewLoggerHandler(settingsLoader),
	)

	if err := gin_middleware_prometheus.SetupPrometheusMetrics(mgr.router, settingsLoader); err != nil {
		return err
	}

	return nil
}

func (mgr *GinManager) setupRoutes(controller icontrollers.CreatureTypes) {
	mgr.router.GET(
		mgr.config.Http.GetCreatureTypesEndpoint,
		handlers.GetCreatureTypesHandler(controller),
	)
	mgr.router.PUT(
		mgr.config.Http.PutCreatureTypesEndpoint,
		handlers.PutCreatureTypesHandler(controller),
	)
	pprof.Register(mgr.router)
}
