package idata_handler

import "creature_types_srv/src/creature_types"

type IDataHandler interface {
	GetCreatureTypes() creature_types.Types
}
