package interfaces

import (
	data_handlers_interfaces "creature_types_srv/src/data_handlers/interfaces"
)

type Router interface {
	Run()
	Stop()
	SetupDataHandler(data_handler data_handlers_interfaces.DataHandler) error
}
