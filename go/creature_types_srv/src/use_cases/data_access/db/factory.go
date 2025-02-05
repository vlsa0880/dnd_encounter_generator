package db

import (
	"fmt"

	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	redisdb "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/data_access/db/redis"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
)

type config struct {
	CreatureTypes struct {
		StorageType string `json:"storage_type"`
	}
}

func New(settingsLoader isettings.SettingsLoader) (interfaces.CreatureTypesRepository, error) {
	config := config{}
	if err := settingsLoader.Load(&config); err != nil {
		return nil, fmt.Errorf("can't load application config: %w", err)
	}
	switch storageType := config.CreatureTypes.StorageType; storageType {
	case "postgres":
		db, err := NewPsql()
		if err != nil {
			return nil, fmt.Errorf("can't construct db '%s':%w", storageType, err)
		}
		return db, nil
	case "redis_postgres":
		db, err := NewPsql()
		if err != nil {
			return nil, fmt.Errorf("can't construct psql: %w", err)
		}
		redisDb, err := redisdb.NewRedisCacheThrough(settingsLoader, db)
		if err != nil {
			return nil, fmt.Errorf("can't construct db '%s':%w", storageType, err)
		}
		return redisDb, nil
	default:
		return nil, fmt.Errorf("no supported db type specified")
	}
}
