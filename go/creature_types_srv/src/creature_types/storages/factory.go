package storages

import (
	"creature_types_srv/src/creature_types"
)

func New(storage_type *string) creature_types.ICreatureTypeManager {
	switch *storage_type {
	case "postgres":
		return &PsqlCreatureTypesManager{}
	}
	return nil
}
