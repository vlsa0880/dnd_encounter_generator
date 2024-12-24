package grpc

import (
	"fmt"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc/services"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type Manager struct {
	router      *GRPCRouter
	dataHandler data_handlers_interfaces.DataHandler
}

func New(settingsLoader settings.SettingsLoader) *Manager {
	var manager Manager
	manager.router = NewServer(settingsLoader)
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

func (manager *Manager) SetupDataHandler(dataHandler data_handlers_interfaces.DataHandler) error {
	if dataHandler == nil {
		return fmt.Errorf("Data handler is nil")
	}
	manager.dataHandler = dataHandler
	manager.registerHandlers()
	return manager.isValid()
}

func (manager *Manager) isValid() error {
	if manager.dataHandler == nil {
		return fmt.Errorf("Data handler is nil")
	} else if manager.router.Server == nil {
		return fmt.Errorf("Grpc server is nil")
	}
	return nil
}

func (manager *Manager) registerHandlers() {
	services.New(manager.router.Server, manager.dataHandler)
}
