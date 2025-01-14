package creature_types_data

import (
	"encoding/json"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
)

type Type struct {
	Name string
}

type Data struct {
	Types []Type
}

func (creatureType *Type) String() string {
	return creatureType.Name
}

func (types Data) String() string {
	json, err := json.Marshal(types)
	if err != nil {
		logger.GetInstance().Warn(
			"Can't convert creature types to string",
		)
		return ""
	}
	return string(json)
}
