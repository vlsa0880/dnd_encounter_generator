package creature_types_storages

import creature_types_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"

type BaseCreatureTypeManager struct {
	creatureTypes creature_types_data.Data
}

func (manager *BaseCreatureTypeManager) GetData() creature_types_data.Data {
	return manager.creatureTypes
}
