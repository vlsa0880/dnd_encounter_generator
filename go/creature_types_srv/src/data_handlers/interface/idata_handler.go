package data_handlers

import "creature_types_srv/src/creature_types"

type IDataHandler interface {
	Init() error
	GetCreatureTypes() creature_types.Types
}
