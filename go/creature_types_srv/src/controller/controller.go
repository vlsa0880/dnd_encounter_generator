package controller

import (
	"creature_types_srv/src/creature_types/storages"
	"creature_types_srv/src/logger"
	"creature_types_srv/src/settings/implementations"
	"fmt"
)

func Run() error {
	settings := implementations.GetInstance().GetSettings()
	logger.GetInstance().Info(settings.GetSlogGroup().String())
	creature_types_mngr := storages.New(&settings.CreatureTypes.StorageType)
	if creature_types_mngr == nil {
		return fmt.Errorf("Can't create creature types mngr")
	}
	var err error
	if err = creature_types_mngr.Init(); err != nil {
		return err
	} else if err = creature_types_mngr.Load(); err != nil {
		return err
	}
	logger.GetInstance().Info(
		fmt.Sprintf("Presented types: %v", creature_types_mngr.GetData()),
	)
	return nil
}
