package creature_types_storages

import creature_types_data "creature_types_srv/src/creature_types/data"

type BaseCreatureTypeManager struct {
	creature_types creature_types_data.Types
}

func (manager *BaseCreatureTypeManager) GetData() creature_types_data.Types {
	return manager.creature_types
}
