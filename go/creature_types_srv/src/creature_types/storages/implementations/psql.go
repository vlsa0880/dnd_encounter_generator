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
	var creatureTypesData []CreatureTypesData
	if err = manager.db.Find(&creatureTypesData).Error; err != nil {
		return err
	}
	manager.fill(creatureTypesData)
	return nil
}

func (manager *PsqlCreatureTypesManager) fill(creatureTypesData []CreatureTypesData) {
	manager.creatureTypes = make(
		creature_types.Types,
		len(creatureTypesData),
	)
	for index, creature_type := range creatureTypesData {
		manager.creatureTypes[index] = creature_type.Type
	}
}
