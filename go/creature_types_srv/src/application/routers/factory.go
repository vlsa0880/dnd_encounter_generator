package routers

import (
	"context"
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/gin"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/grpc"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka"

	"go.uber.org/zap"
)

type Config struct {
	Router struct {
		Types []string
	}
}

func New(ctx context.Context, settingsLoader settings.SettingsLoader, creatureTypesDB dbinterfaces.CreatureTypes) ([]interfaces.Router, error) {
	if settingsLoader == nil {
		return nil, fmt.Errorf("bad settings loader")
	}

	config := Config{}
	if err := settingsLoader.Load(&config); err != nil {
		return nil, fmt.Errorf("can't load routers factory settings: %w", err)
	}

	routers := make([]interfaces.Router, 0, len(config.Router.Types))
	for _, routerName := range config.Router.Types {
		switch routerName {
		case "rest":
			ginRouter, err := gin.New(settingsLoader, creatureTypesDB)
			if err != nil {
				return nil, fmt.Errorf("can't create gin router: %w", err)
			}
			routers = append(routers, ginRouter)
		case "grpc":
			grpcManager, err := grpc.New(settingsLoader, creatureTypesDB)
			if err != nil {
				return nil, fmt.Errorf("can't create grpc router: %w", err)
			}
			routers = append(routers, grpcManager)
		case "kafka":
			kafkaManager, err := kafka.New(ctx, settingsLoader, creatureTypesDB)
			if err != nil {
				return nil, fmt.Errorf("can't create kafka router: %w", err)
			}
			routers = append(routers, kafkaManager)
		default:
			logger.GetInstance().Warn(
				"unknown router type",
				zap.String("router_name", routerName),
			)
		}
	}
	if len(routers) <= 0 {
		return nil, fmt.Errorf("no valid router names found")
	}
	return routers, nil
}
