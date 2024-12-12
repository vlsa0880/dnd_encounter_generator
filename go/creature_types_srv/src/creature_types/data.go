package creature_types

import (
	logger "creature_types_srv/src/logger/zap"
	"encoding/json"
)

type Type struct {
	RU string
}

type Types []Type

func (creature_type *Type) String() string {
	return creature_type.RU
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
