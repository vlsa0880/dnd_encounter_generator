package storages

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgtype"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	settings "creature_types_srv/src/settings/manager"
	utils_db "creature_types_srv/src/utils/db"
)

type PsqlJSONB struct {
	BaseSettingsManager
	db *gorm.DB
}

type Model struct {
	gorm.Model
	Data pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
}

func (settings_manager *PsqlJSONB) Init() error {
	dsn := utils_db.DsnByEnv()
	var err error
	fmt.Println(*dsn)
	settings_manager.db, err = gorm.Open(postgres.Open(*dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	return nil
}

func (settings_manager *PsqlJSONB) Load() error {
	settings_manager.settings = settings.NewData()
	model := Model{}
	var err error
	if err = model.Data.Set(settings_manager.settings); err != nil {
		return err
	}
	if err = settings_manager.db.AutoMigrate(&model); err != nil {
		return err
	}
	err = settings_manager.db.First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = settings_manager.db.Model(&model).FirstOrCreate(&model).Error
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if err = settings_manager.db.Find(&model).Error; err != nil {
			return err
		}
	}

	if err = jsonbTransform(&model.Data, &settings_manager.settings); err != nil {
		return err
	}
	return nil
}

func jsonbTransform(jsonb *pgtype.JSONB, res interface{}) error {
	values, err := json.Marshal(jsonb)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(values, res); err != nil {
		return err
	}
	return nil
}
