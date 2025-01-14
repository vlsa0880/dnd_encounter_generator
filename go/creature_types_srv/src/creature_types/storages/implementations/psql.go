package creature_types_storages

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	global_logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"moul.io/zapgorm2"

	creature_types "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	creature_types_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	db_utils "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/utils/db"
)

type PsqlCreatureTypesManager struct {
	db *gorm.DB
}

type CreatureTypesData struct {
	creature_types.Type
	gorm.Model
}

func (manager *PsqlCreatureTypesManager) Init() error {
	dsn := db_utils.DsnByEnv()
	var err error
	logger := zapgorm2.New(global_logger.GetInstance())
	manager.db, err = gorm.Open(
		postgres.Open(*dsn),
		&gorm.Config{
			Logger: logger,
		})
	if err != nil {
		return err
	}
	return nil
}

func (manager *PsqlCreatureTypesManager) GetData() (*creature_types_data.Data, error) {
	var err error
	if err = manager.db.AutoMigrate(&CreatureTypesData{}); err != nil {
		return nil, err
	}
	var creatureTypesData []CreatureTypesData
	if err = manager.db.Find(&creatureTypesData).Error; err != nil {
		return nil, err
	}
	res := creature_types_data.Data{}
	for _, creature_type := range creatureTypesData {
		res.Types = append(res.Types, creature_types_data.Type{Name: creature_type.Name})
	}
	return &res, nil
}

func (manager *PsqlCreatureTypesManager) Upload(types *creature_types.Data) error {
	return manager.db.Transaction(func(tx *gorm.DB) error {
		if err := manager.db.Delete(&CreatureTypesData{}, "*").Error; err != nil {
			return fmt.Errorf("can't delete all records from creature type table")
		}
		if err := manager.db.Create(&types).Error; err != nil {
			return fmt.Errorf("can't upload creature types: %w", err)
		}
		return nil
	})
}
