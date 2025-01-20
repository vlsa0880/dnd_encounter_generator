package grpc

import (
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers/grpc/services"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
	"go.uber.org/zap"
)

type Manager struct {
	router          *GRPCRouter
	creatureTypesDB dbinterfaces.CreatureTypes
}

func New(settingsLoader settings.SettingsLoader, creatureTypesDB dbinterfaces.CreatureTypes) *Manager {
	var manager Manager
	manager.router = NewServer(settingsLoader)
	if err := manager.setupDataHandler(creatureTypesDB); err != nil {
		logger.GetInstance().Error(
			"can't setup data handler",
			zap.Error(err),
		)
		return nil
	}
	return &manager
}

func (manager *Manager) Run() {
	if err := manager.isValid(); err != nil {
		err_msg := fmt.Sprintf("Can't run grpc server: %s", err)
		panic(err_msg)
	}
	manager.router.Run()
}

func (manager *Manager) Stop() {
	if err := manager.isValid(); err != nil {
		err_msg := fmt.Sprintf("Can't stop grpc server: %s", err)
		panic(err_msg)
	}
	manager.router.Stop()
}

func (manager *Manager) setupDataHandler(creatureTypesDB dbinterfaces.CreatureTypes) error {
	if creatureTypesDB == nil {
		return fmt.Errorf("Data handler is nil")
	}
	manager.creatureTypesDB = creatureTypesDB
	manager.registerHandlers()
	return manager.isValid()
}

func (manager *Manager) isValid() error {
	if manager.creatureTypesDB == nil {
		return fmt.Errorf("Data handler is nil")
	} else if manager.router.Server == nil {
		return fmt.Errorf("Grpc server is nil")
	}
	return nil
}

func (manager *Manager) registerHandlers() {
	services.New(manager.router.Server, manager.creatureTypesDB)
}
