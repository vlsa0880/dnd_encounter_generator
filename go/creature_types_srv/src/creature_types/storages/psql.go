package storages

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"creature_types_srv/src/creature_types"
	db_utils "creature_types_srv/src/utils/db"
)

// TODO: for shadow table updates
// const c_tmp_table_name string = "tmp_creature_types"

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
		Logger: logger.Default.LogMode(logger.Info),
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

func (manager *PsqlCreatureTypesManager) fill(data []CreatureTypesData) {
	manager.creature_types = make(
		creature_types.Types,
		len(data),
	)
	for index, creature_type := range data {
		manager.creature_types[index] = creature_type.Type
	}
}
