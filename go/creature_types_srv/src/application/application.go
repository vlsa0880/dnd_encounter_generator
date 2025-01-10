package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	data_handlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/handlers"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	routes_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type RouterProcess func(router routes_interfaces.Router)

type Application struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	routers   []routes_interfaces.Router
}

func NewApplication() *Application {
	loader := settings.New()
	if loader == nil {
		return nil
	}
	var settingsLoader isettings.SettingsLoader = loader

	logger.InitLogger(settingsLoader)

	handler := data_handlers.New(settingsLoader)
	if handler == nil {
		logger.GetInstance().Error("can't create data handler")
		return nil
	}

	application := Application{}
	application.ctx, application.ctxCancel = context.WithCancel(context.Background())
	application.routers = routers.New(application.ctx, settingsLoader, handler)
	if len(application.routers) < 1 {
		panic("no routers created - check env")
	}

	return &application
}

func (application *Application) Run() {
	defer application.ctxCancel()
	application.forEach(func(router routes_interfaces.Router) { go router.Run() })
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.forEach(func(router routes_interfaces.Router) { go router.Stop() })
}

func (application *Application) forEach(processer RouterProcess) {
	for _, router := range application.routers {
		if router == nil {
			panic("some of application routers is nil")
		}
		processer(router)
	}
}
