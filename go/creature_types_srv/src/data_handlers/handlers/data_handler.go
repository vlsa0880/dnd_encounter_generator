package data_handlers

import (
	"context"

	creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	storages "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/storages/implementations"
	creature_types_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/storages/interfaces"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"

	"go.uber.org/zap"
)

type Config struct {
	CreatureTypes struct {
		GetEndpoint string `json:"get_endpoint"`
		StorageType string `json:"storage_type"`
	}
}

type DataHandler struct {
	creatureTypesManager creature_types_interfaces.CreatureTypeManager
}

func New(settings_loader settings.SettingsLoader) *DataHandler {
	config := Config{}
	var err error
	if err = settings_loader.Load(&config); err != nil {
		logger.GetInstance().Error("Can't load config")
		return nil
	}
	handler := DataHandler{}
	handler.creatureTypesManager = storages.New(&config.CreatureTypes.StorageType)
	if handler.creatureTypesManager == nil {
		logger.GetInstance().Error("Can't create creature types mngr")
		return nil
	}
	if err = handler.creatureTypesManager.Init(); err != nil {
		logger.GetInstance().Error(
			"Can't init creature types mngr",
			zap.String("msg", err.Error()),
		)
		return nil
	} else if err = handler.creatureTypesManager.Load(); err != nil {
		logger.GetInstance().Error(
			"Can't load types from creature types mngr",
			zap.String("msg", err.Error()),
		)
		return nil
	}
	return &handler
}

func (handler *DataHandler) GetCreatureTypes(ctx context.Context) creature_types.Types {
	return handler.creatureTypesManager.GetData()
}
