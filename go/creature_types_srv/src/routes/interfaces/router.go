package interfaces

import (
	"creature_types_srv/src/data_handlers/interfaces"
)

type Router interface {
	Run()
	Stop()
	SetupDataHandler(data_handler interfaces.DataHandler) error
}
