package storages

import (
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/manager"
)

type BaseSettingsManager struct {
	settings *settings.Data
}

func (base_settings *BaseSettingsManager) GetSettings() *settings.Data {
	return base_settings.settings
}

func (base_settings *BaseSettingsManager) Load() error {
	base_settings.settings = settings.NewData()
	return nil
}
