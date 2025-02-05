package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

type RedisCacheThrough struct {
	repository interfaces.CreatureTypesRepository
	client     *redis.Client
	config     config
	updater    redisUpdater
}

type config struct {
	Redis struct {
		Password         string
		CreatureTypesKey string
		Connect          struct {
			Address string
			Timeout time.Duration
		}
	}
}

func NewRedisCacheThrough(settingsLoader isettings.SettingsLoader, db interfaces.CreatureTypesRepository) (*RedisCacheThrough, error) {
	if settingsLoader == nil {
		return nil, fmt.Errorf("bad settings loader")
	}

	cache := RedisCacheThrough{}

	if err := settingsLoader.Load(&cache.config); err != nil {
		return nil, fmt.Errorf("can't load config: %w", err)
	}

	cache.client = redis.NewClient(
		&redis.Options{
			Addr:     cache.config.Redis.Connect.Address,
			Password: cache.config.Redis.Password,
		})

	timeoutCtx, cancel := context.WithTimeout(context.Background(), cache.config.Redis.Connect.Timeout)
	defer cancel()
	if _, err := cache.client.Ping(timeoutCtx).Result(); err != nil {
		return nil, fmt.Errorf("can't connect to redis: %w", err)
	}

	var err error
	if cache.updater, err = newUpdater(settingsLoader, db, cache.client); err != nil {
		return nil, fmt.Errorf("can't create updater: %w", err)
	}
	if err := cache.updater.updateCacheWithDbData(timeoutCtx); err != nil {
		return nil, fmt.Errorf("can't update cache: %w", err)
	}

	return &cache, nil
}

func (cache *RedisCacheThrough) GetCreatureTypes(ctx context.Context) (entities.CreatureTypes, error) {
	if !cache.updater.ready() {
		logger.GetInstance().Debug("updater isn't ready, return data from db")
		return cache.repository.GetCreatureTypes(ctx)
	}

	data, err := cache.client.Get(ctx, cache.config.Redis.CreatureTypesKey).Bytes()
	if err != nil {
		logger.GetInstance().Error(
			"can't get data from cache",
			zap.Error(err),
		)
		return cache.repository.GetCreatureTypes(ctx)
	}

	var res entities.CreatureTypes
	err = msgpack.Unmarshal(data, &res)
	if err != nil {
		logger.GetInstance().Error(
			"can't deserialize data from cache",
			zap.Error(err),
		)
		return cache.repository.GetCreatureTypes(ctx)
	}
	logger.GetInstance().Debug("successfully load data from cache")
	return res, nil
}

func (cache *RedisCacheThrough) UploadCreatureTypes(ctx context.Context, types *entities.CreatureTypes) error {
	var err error
	if err = cache.repository.UploadCreatureTypes(ctx, types); err != nil {
		return fmt.Errorf("can't upload data to db: %w", err)
	}
	if err = cache.updater.update(ctx, types); err != nil {
		return fmt.Errorf("can't upload data to cache: %w", err)
	}
	return nil
}
