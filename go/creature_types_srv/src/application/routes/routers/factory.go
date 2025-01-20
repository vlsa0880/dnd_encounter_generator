package routers

import (
	"context"
	"fmt"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers/gin"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers/grpc"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers/kafka"

	"go.uber.org/zap"
)

type Config struct {
	Router struct {
		Types []string
	}
}

func New(ctx context.Context, settingsLoader settings.SettingsLoader, creatureTypesDB dbinterfaces.CreatureTypes) []interfaces.Router {
	if settingsLoader == nil {
		return []interfaces.Router{}
	}
	config := Config{}
	if err := settingsLoader.Load(&config); err != nil {
		err_msg := fmt.Errorf("can't load routers factory settings: %s", err)
		panic(err_msg)
	}

	routers := []interfaces.Router{}
	for _, router_name := range config.Router.Types {
		switch router_name {
		case "rest":
			gin_router := gin.New(settingsLoader, creatureTypesDB)
			if gin_router == nil {
				panic("can't create gin router")
			}
			routers = append(routers, gin_router)
		case "grpc":
			grpc_router := grpc.New(settingsLoader, creatureTypesDB)
			if grpc_router == nil {
				panic("can't create grpc router")
			}
			routers = append(routers, grpc_router)
		case "kafka":
			kafkaRouter := kafka.New(ctx, settingsLoader, creatureTypesDB)
			if kafkaRouter == nil {
				panic("can't create kafka router")
			}
			routers = append(routers, kafkaRouter)
		default:
			logger.GetInstance().Warn(
				"unknown router type",
				zap.String("router_name", router_name),
			)
		}
	}
	return routers
}
