package storages

import (
	"creature_types_srv/src/creature_types"
)

type BaseCreatureTypeManager struct {
	creature_types creature_types.Types
}

func (manager *BaseCreatureTypeManager) GetData() creature_types.Types {
	return manager.creature_types
}
