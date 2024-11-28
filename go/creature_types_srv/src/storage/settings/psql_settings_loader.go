package settings

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgtype"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PsqlSettingsManager struct {
	BaseSettingsManager
	db *gorm.DB
}

type CreatureTypesSettings struct {
	gorm.Model
	Data pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null"`
}

func (settings_manager *PsqlSettingsManager) Init() error {
	var err error
	if err = settings_manager.openDB(); err != nil {
		return err
	}
	settings_data := newSettingsData()
	creature_settings := CreatureTypesSettings{}
	creature_settings.Data.Set(settings_data)
	if err = settings_manager.db.AutoMigrate(&creature_settings); err != nil {
		return err
	}
	var creature_settings_exists bool
	err = settings_manager.db.Model(&creature_settings).
		Select("count(*) > 0").
		Find(&creature_settings_exists).
		Error
	if err != nil {
		return err
	}

	if !creature_settings_exists {
		err = settings_manager.db.Model(&creature_settings).FirstOrCreate(&creature_settings).Error
		if err != nil {
			return err
		}
	} else {
		if err = settings_manager.db.Find(&creature_settings).Error; err != nil {
			return err
		}
	}

	if err = jsonbTransform(&creature_settings.Data, &settings_manager.settings); err != nil {
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

func (settings_manager *PsqlSettingsManager) openDB() error {
	storage_address := os.Getenv(c_storage_address_env)
	if len(storage_address) == 0 {
		storage_address = c_default_storage_address
	}
	storage_name := os.Getenv(c_storage_name_env)
	if len(storage_name) == 0 {
		storage_name = c_default_storage_name
	}

	dsn := fmt.Sprintf("host=%v user=postgres dbname=%v password=test port=5432", storage_address, storage_name)
	var err error
	settings_manager.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	return nil
}
