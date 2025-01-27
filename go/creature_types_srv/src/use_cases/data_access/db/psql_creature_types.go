package db

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	global_logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"moul.io/zapgorm2"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	db_utils "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/utils/db"
)

type PsqlCreatureTypesController struct {
	db *gorm.DB
}

type dbCreatureTypesData struct {
	gorm.Model
	entities.CreatureType
}

func NewPsql() (*PsqlCreatureTypesController, error) {
	dsn := db_utils.DsnByEnv()
	var err error
	logger := zapgorm2.New(global_logger.GetInstance())
	controller := PsqlCreatureTypesController{}
	controller.db, err = gorm.Open(
		postgres.Open(*dsn),
		&gorm.Config{
			Logger: logger,
		})
	if err != nil {
		return nil, err
	}
	if err = controller.db.AutoMigrate(&dbCreatureTypesData{}); err != nil {
		return nil, err
	}
	return &controller, nil
}

func (controller *PsqlCreatureTypesController) GetCreatureTypes(ctx context.Context) (entities.CreatureTypes, error) {
	var err error
	var creatureTypesData []dbCreatureTypesData
	if err = controller.db.Find(&creatureTypesData).Error; err != nil {
		return entities.CreatureTypes{}, err
	}
	res := entities.CreatureTypes{
		Types: make([]entities.CreatureType, 0, len(creatureTypesData)),
	}
	for _, dbCreatureTypeData := range creatureTypesData {
		global_logger.GetInstance().Warn(
			"get creature type data",
			zap.String("data", fmt.Sprintf("%v", dbCreatureTypeData)),
		)
		creatureType := entities.CreatureType{
			ID:   dbCreatureTypeData.Model.ID,
			Name: dbCreatureTypeData.CreatureType.Name,
		}
		global_logger.GetInstance().Warn(
			"creatureType",
			zap.String("creatureType", fmt.Sprintf("%v", creatureType)),
		)
		res.Types = append(
			res.Types,
			creatureType,
		)
	}
	global_logger.GetInstance().Warn(
		"res",
		zap.String("res", fmt.Sprintf("%v", res)),
	)
	return res, nil
}

func (controller *PsqlCreatureTypesController) UploadCreatureTypes(ctx context.Context, types *entities.CreatureTypes) error {
	return controller.db.Transaction(func(tx *gorm.DB) error {
		if err := controller.db.Where("deleted_at is NULL").Delete(&dbCreatureTypesData{}).Error; err != nil {
			return fmt.Errorf("can't delete all records from creature type table")
		}
		sqlData := make([]dbCreatureTypesData, 0, len(types.Types))
		for _, dbCreatureTypeData := range types.Types {
			row := dbCreatureTypesData{
				CreatureType: entities.CreatureType{
					Name: dbCreatureTypeData.Name,
				},
			}
			sqlData = append(sqlData, row)
		}
		if err := controller.db.Create(&sqlData).Error; err != nil {
			return fmt.Errorf("can't upload creature types: %w", err)
		}
		return nil
	})
}
