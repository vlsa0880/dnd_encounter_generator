package db

import (
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

func New(storageType *string) (interfaces.CreatureTypesRepository, error) {
	switch *storageType {
	case "postgres":
		controller, err := NewPsql()
		if err != nil {
			return nil, fmt.Errorf("can't build creature types db '%s':%w", *storageType, err)
		}
		return controller, nil
	default:
		return nil, fmt.Errorf("no supported db type specified")
	}
}
