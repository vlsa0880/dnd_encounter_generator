package irouter

import (
	"creature_types_srv/src/data_handlers/idata_handler"
)

type IRouter interface {
	Init() error
	Run()
	SetupDataHandler(data_handler idata_handler.IDataHandler) error
}
