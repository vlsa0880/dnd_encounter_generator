package grpc

import (
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/grpc/services"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type Manager struct {
	router          *GRPCRouter
	creatureTypesDB dbinterfaces.CreatureTypes
}

func New(settingsLoader settings.SettingsLoader, creatureTypesDB dbinterfaces.CreatureTypes) (*Manager, error) {
	var manager Manager
	var err error
	manager.router, err = NewServer(settingsLoader)
	if err != nil {
		return nil, fmt.Errorf("can't construct grpc server: %w", err)
	}
	if err := manager.setupDataHandler(creatureTypesDB); err != nil {
		return nil, fmt.Errorf("can't setup data handler: %w", err)
	}
	return &manager, nil
}

func (manager *Manager) Run() error {
	if err := manager.isValid(); err != nil {
		return fmt.Errorf("Can't run grpc server: %w", err)
	}
	return manager.router.Run()
}

func (manager *Manager) Stop() error {
	if err := manager.isValid(); err != nil {
		err_msg := fmt.Sprintf("Can't stop grpc server: %s", err)
		panic(err_msg)
	}
	return manager.router.Stop()
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
	services.Run(manager.router.Server, manager.creatureTypesDB)
}
