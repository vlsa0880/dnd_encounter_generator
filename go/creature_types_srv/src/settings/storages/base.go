package storages

import (
	settings "creature_types_srv/src/settings/manager"
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
