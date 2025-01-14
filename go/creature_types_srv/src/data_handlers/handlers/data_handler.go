package data_handlers

import (
	"context"
	"fmt"

	creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	storages "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/storages/implementations"
	creature_types_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/storages/interfaces"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"

	"go.uber.org/zap"
)

type Config struct {
	CreatureTypes struct {
		StorageType string `json:"storage_type"`
	}
}

type DataHandler struct {
	creatureTypesManager creature_types_interfaces.CreatureTypeManager
}

func New(settingsLoader settings.SettingsLoader) *DataHandler {
	config := Config{}
	var err error
	if err = settingsLoader.Load(&config); err != nil {
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
			zap.Error(err),
		)
		return nil
	}
	return &handler
}

func (handler *DataHandler) GetCreatureTypes(ctx context.Context) (*creature_types.Data, error) {
	return handler.creatureTypesManager.GetData()
}

func (handler *DataHandler) PutCreatureTypes(ctx context.Context, types *creature_types.Data) error {
	return fmt.Errorf("not implemented")
}
