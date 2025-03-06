package ginrouter

import (
	"fmt"

	handlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/gin/handlers"
	gin_middleware "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/gin/middleware"
	gin_middleware_prometheus "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/gin/middleware/prometheus"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"

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
	router          *gin.Engine
	creatureTypesDB dbinterfaces.CreatureTypesRepository
	config          Config
}

func New(settingsLoader settings.SettingsLoader, creatureTypesDB dbinterfaces.CreatureTypesRepository) (*GinManager, error) {
	mgr := GinManager{}
	if err := settingsLoader.Load(&mgr.config); err != nil {
		return nil, fmt.Errorf("Error loading gin router config: %w", err)
	}
	if creatureTypesDB == nil {
		return nil, fmt.Errorf("bad data handler")
	}
	mgr.creatureTypesDB = creatureTypesDB

	mgr.router = gin.New()
	if err := mgr.setupMiddleware(settingsLoader); err != nil {
		return nil, fmt.Errorf("can't setup middleware: %w", err)
	}

	mgr.setupRoutes(creatureTypesDB)
	return &mgr, nil
}

func (mgr *GinManager) Run() error {
	if mgr.router == nil {
		return fmt.Errorf("Bad router")
	}
	if mgr.creatureTypesDB == nil {
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

func (mgr *GinManager) setupRoutes(creatureTypesDB dbinterfaces.CreatureTypesRepository) {
	mgr.router.GET(
		mgr.config.Http.GetCreatureTypesEndpoint,
		handlers.GetCreatureTypesHandler(creatureTypesDB),
	)
	mgr.router.PUT(
		mgr.config.Http.PutCreatureTypesEndpoint,
		handlers.PutCreatureTypesHandler(creatureTypesDB),
	)
	pprof.Register(mgr.router)
}
