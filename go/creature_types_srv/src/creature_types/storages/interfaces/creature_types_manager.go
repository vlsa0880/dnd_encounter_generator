package creature_types_interfaces

import (
	creature_types_data "creature_types_srv/src/creature_types/data"
)

type CreatureTypeManager interface {
	Init() error
	Load() error
	GetData() creature_types_data.Types
}
