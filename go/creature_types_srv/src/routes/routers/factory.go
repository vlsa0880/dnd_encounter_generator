package routers

import (
	logger "creature_types_srv/src/logger/zap"
	"creature_types_srv/src/routes/interfaces"
	settings "creature_types_srv/src/settings/loader"
	"fmt"

	"creature_types_srv/src/routes/routers/gin"
	"creature_types_srv/src/routes/routers/grpc"

	"go.uber.org/zap"
)

type Config struct {
	Router struct {
		Types []string
	}
}

func New(settings_loader settings.ISettingsLoader) []interfaces.Router {
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
