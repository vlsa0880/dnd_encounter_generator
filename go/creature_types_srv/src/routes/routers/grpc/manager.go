package grpc

import (
	"fmt"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc/services"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"go.uber.org/zap"
)

type Manager struct {
	router      *GRPCRouter
	dataHandler data_handlers_interfaces.DataHandler
}

func New(settingsLoader settings.SettingsLoader, dataHandler data_handlers_interfaces.DataHandler) *Manager {
	var manager Manager
	manager.router = NewServer(settingsLoader)
	if err := manager.setupDataHandler(dataHandler); err != nil {
		logger.GetInstance().Error(
			"can't setup data handler",
			zap.String("msg", err.Error()),
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

func (manager *Manager) setupDataHandler(dataHandler data_handlers_interfaces.DataHandler) error {
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
