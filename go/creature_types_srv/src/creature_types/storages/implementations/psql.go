package creature_types_storages

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	db_utils "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/utils/db"
)

type PsqlCreatureTypesManager struct {
	BaseCreatureTypeManager
	db *gorm.DB
}

type CreatureTypesData struct {
	creature_types.Type
	gorm.Model
}

func (manager *PsqlCreatureTypesManager) Init() error {
	dsn := db_utils.DsnByEnv()
	var err error
	manager.db, err = gorm.Open(postgres.Open(*dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	return nil
}

func (manager *PsqlCreatureTypesManager) Load() error {
	var err error
	if err = manager.db.AutoMigrate(&CreatureTypesData{}); err != nil {
		return err
	}
	var creatures_types_data []CreatureTypesData
	if err = manager.db.Find(&creatures_types_data).Error; err != nil {
		return err
	}
	manager.fill(creatures_types_data)
	return nil
}

func (manager *PsqlCreatureTypesManager) fill(creature_types_data []CreatureTypesData) {
	manager.creature_types = make(
		creature_types.Types,
		len(creature_types_data),
	)
	for index, creature_type := range creature_types_data {
		manager.creature_types[index] = creature_type.Type
	}
}
