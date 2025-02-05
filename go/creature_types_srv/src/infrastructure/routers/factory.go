package routers

import (
	"context"
	"fmt"

	icontrollers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	ginrouter "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/gin"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/grpc"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/kafka"

	"go.uber.org/zap"
)

type Config struct {
	Router struct {
		Types []string
	}
}

func New(ctx context.Context, settingsLoader settings.SettingsLoader, controller icontrollers.CreatureTypes) ([]interfaces.Router, error) {
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
			ginRouter, err := ginrouter.New(settingsLoader, controller)
			if err != nil {
				return nil, fmt.Errorf("can't create gin router: %w", err)
			}
			routers = append(routers, ginRouter)
		case "grpc":
			grpcManager, err := grpc.New(settingsLoader, controller)
			if err != nil {
				return nil, fmt.Errorf("can't create grpc router: %w", err)
			}
			routers = append(routers, grpcManager)
		case "kafka":
			kafkaManager, err := kafka.New(ctx, settingsLoader, controller)
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
