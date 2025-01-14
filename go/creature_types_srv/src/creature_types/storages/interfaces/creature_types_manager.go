package creature_types_interfaces

import (
	creature_types_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
)

type CreatureTypeManager interface {
	Init() error
	GetData() (*creature_types_data.Data, error)
	Upload(types *creature_types_data.Data) error
}
