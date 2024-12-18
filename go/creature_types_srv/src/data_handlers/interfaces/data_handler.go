package interfaces

import "creature_types_srv/src/creature_types"

type DataHandler interface {
	GetCreatureTypes() creature_types.Types
}
