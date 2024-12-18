package handlers

import (
	"creature_types_srv/src/creature_types"
	"creature_types_srv/src/creature_types/storages"
	settings "creature_types_srv/src/settings/loader"

	logger "creature_types_srv/src/logger/zap"

	"go.uber.org/zap"
)

type Config struct {
	CreatureTypes struct {
		GetEndpoint string `json:"get_endpoint"`
		StorageType string `json:"storage_type"`
	}
}

type DataHandler struct {
	creature_types_mngr creature_types.CreatureTypeManager
}

func New(settings_loader settings.ISettingsLoader) *DataHandler {
	config := Config{}
	var err error
	if err = settings_loader.Load(&config); err != nil {
		logger.GetInstance().Error("Can't load config")
		return nil
	}
	handler := DataHandler{}
	handler.creature_types_mngr = storages.New(&config.CreatureTypes.StorageType)
	if handler.creature_types_mngr == nil {
		logger.GetInstance().Error("Can't create creature types mngr")
		return nil
	}
	if err = handler.creature_types_mngr.Init(); err != nil {
		logger.GetInstance().Error(
			"Can't init creature types mngr",
			zap.String("msg", err.Error()),
		)
		return nil
	} else if err = handler.creature_types_mngr.Load(); err != nil {
		logger.GetInstance().Error(
			"Can't load types from creature types mngr",
			zap.String("msg", err.Error()),
		)
		return nil
	}
	return &handler
}

func (handler *DataHandler) GetCreatureTypes() creature_types.Types {
	return handler.creature_types_mngr.GetData()
}
