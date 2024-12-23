package storages

import (
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/utils/env"
)

func CreateSettingsManager() interface{} {
	settings_storage_type := env.Get("SETTINGS_STORAGE_TYPE", "env")
	var imanager interface{}
	switch *settings_storage_type {
	case "env":
		imanager = &Env{}
	case "postgres_jsonb":
		imanager = &PsqlJSONB{}
	default:
		imanager = nil
	}
	return imanager
}
