package interfaces

import (
	data_handlers_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
)

type Router interface {
	Run()
	Stop()
	Init(data_handler data_handlers_interfaces.DataHandler) error
}
