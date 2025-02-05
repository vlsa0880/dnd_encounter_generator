package redisdb

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	readyStatus uint32 = iota
	updatingStatus
)

type redisUpdater struct {
	client       *redis.Client
	repository   interfaces.CreatureTypesRepository
	config       updaterConfig
	updateStatus *atomic.Uint32
}

type updaterConfig struct {
	Redis struct {
		CreatureTypesKey string
		Update           struct {
			Timeout time.Duration
			Sleep   time.Duration
		}
	}
}

func newUpdater(settingsLoader isettings.SettingsLoader, db interfaces.CreatureTypesRepository, client *redis.Client) (redisUpdater, error) {
	if settingsLoader == nil {
		return redisUpdater{}, fmt.Errorf("bad settings loader")
	}
	if db == nil {
		return redisUpdater{}, fmt.Errorf("bad db interface")
	}
	if client == nil {
		return redisUpdater{}, fmt.Errorf("bad redis client")
	}

	updater := redisUpdater{
		repository:   db,
		updateStatus: &atomic.Uint32{},
		client:       client,
	}
	if err := settingsLoader.Load(&updater.config); err != nil {
		return redisUpdater{}, fmt.Errorf("can't load config: %w", err)
	}
	return updater, nil
}

func (updater *redisUpdater) updateCacheWithDbData(ctx context.Context) error {
	types, err := updater.repository.GetCreatureTypes(ctx)
	if err != nil {
		return fmt.Errorf("can't update cache with db value: %w", err)
	}
	return updater.update(ctx, &types)
}

func (updater *redisUpdater) update(ctx context.Context, types *entities.CreatureTypes) error {
	cacheData, err := msgpack.Marshal(types)
	if err != nil {
		updater.runBackgroundUpdater()
		return fmt.Errorf("can't deserialize data from cache: %w", err)
	}
	if err = updater.client.Set(ctx, updater.config.Redis.CreatureTypesKey, cacheData, 0).Err(); err != nil {
		updater.runBackgroundUpdater()
		return fmt.Errorf("%w", err)
	}
	updater.updateStatus.Store(readyStatus)
	return nil
}

func (updater *redisUpdater) runBackgroundUpdater() {
	if updater.updateStatus.Load() == updatingStatus {
		return
	}

	logger.GetInstance().Debug("background updater started")
	updater.updateStatus.Store(updatingStatus)
	go func() {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), updater.config.Redis.Update.Timeout)
		defer cancel()
		for updater.updateCacheWithDbData(ctxTimeout) != nil {
			time.Sleep(updater.config.Redis.Update.Sleep)
		}
		logger.GetInstance().Debug("background updater finished")
	}()
}

func (updater *redisUpdater) ready() bool {
	return updater.updateStatus.Load() == readyStatus
}
