package creature_types_data

import (
	"encoding/json"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
)

type Type struct {
	Name string
}

type Types []Type

func (creature_type *Type) String() string {
	return creature_type.Name
}

func (types Types) String() string {
	json, err := json.Marshal(types)
	if err != nil {
		logger.GetInstance().Warn(
			"Can't convert creature types to string",
		)
		return ""
	}
	return string(json)
}
