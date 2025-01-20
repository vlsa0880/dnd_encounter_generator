package db

import "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"

type BaseCreatureTypeManager struct {
	creatureTypes entities.CreatureTypes
}

func (manager *BaseCreatureTypeManager) GetData() entities.CreatureTypes {
	return manager.creatureTypes
}
