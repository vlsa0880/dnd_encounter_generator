package routers

import (
	"fmt"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/interfaces"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/gin"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/grpc"

	"go.uber.org/zap"
)

type Config struct {
	Router struct {
		Types []string
	}
}

func New(settings_loader settings.SettingsLoader) []interfaces.Router {
	if settings_loader == nil {
		return []interfaces.Router{}
	}
	config := Config{}
	if err := settings_loader.Load(&config); err != nil {
		err_msg := fmt.Errorf("can't load routers factory settings: %s", err)
		panic(err_msg)
	}

	routers := []interfaces.Router{}
	for _, router_name := range config.Router.Types {
		switch router_name {
		case "rest":
			gin_router := gin.New(settings_loader)
			if gin_router == nil {
				panic("can't create gin router")
			}
			routers = append(routers, gin_router)
		case "grpc":
			grpc_router := grpc.New(settings_loader)
			if grpc_router == nil {
				panic("can't create grpc router")
			}
			routers = append(routers, grpc_router)
		default:
			logger.GetInstance().Warn(
				"unknown router type",
				zap.String("router_name", router_name),
			)
		}
	}
	return routers
}
