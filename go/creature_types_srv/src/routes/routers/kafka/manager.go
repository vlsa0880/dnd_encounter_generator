package kafka

import (
	"fmt"

	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
)

type Manager struct {
}

func New() *Manager {
	return &Manager{}
}

func (manager *Manager) Run() {
	panic("not implemented")
}

func (manager *Manager) Stop() {
	panic("not implemented")
}

func (manager *Manager) SetupDataHandler(data_handler data_handlers_interfaces.DataHandler) error {
	return fmt.Errorf("not implemented")
}
