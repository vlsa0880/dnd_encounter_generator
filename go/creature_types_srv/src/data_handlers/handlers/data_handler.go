package handlers

import (
	"creature_types_srv/src/creature_types"
	"creature_types_srv/src/creature_types/storages"
	"creature_types_srv/src/settings/implementations"
	"fmt"
)

type DataHandler struct {
	creature_types_mngr creature_types.ICreatureTypeManager
}

func (handler *DataHandler) Init() error {
	settings := implementations.GetInstance().GetSettings()
	handler.creature_types_mngr = storages.New(&settings.CreatureTypes.StorageType)
	if handler.creature_types_mngr == nil {
		return fmt.Errorf("Can't create creature types mngr")
	}
	var err error
	if err = handler.creature_types_mngr.Init(); err != nil {
		return err
	} else if err = handler.creature_types_mngr.Load(); err != nil {
		return err
	}
	return nil
}

func (handler *DataHandler) GetCreatureTypes() creature_types.Types {
	return handler.creature_types_mngr.GetData()
}
