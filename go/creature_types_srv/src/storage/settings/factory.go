package settings

import "os"

func CreateSettingsManager() ISettingsManager {
	settings_storage_type := os.Getenv(c_settings_storage_type_env)
	if len(settings_storage_type) == 0 {
		settings_storage_type = c_default_settings_storage_type
	}
	var imanager ISettingsManager
	switch settings_storage_type {
	case "postgres":
		imanager = &PsqlSettingsManager{}
	default:
		imanager = nil
	}
	return imanager
}
