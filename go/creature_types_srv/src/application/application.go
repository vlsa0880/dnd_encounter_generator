package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	irouter "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routes/routers"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/data_access/db"
	"go.uber.org/zap"
)

type RouterProcess func(router irouter.Router)

type config struct {
	CreatureTypes struct {
		StorageType string `json:"storage_type"`
	}
}

type Application struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	routers   []irouter.Router
}

func NewApplication() *Application {
	loader := settings.New()
	if loader == nil {
		return nil
	}
	var settingsLoader isettings.SettingsLoader = loader

	logger.InitLogger(settingsLoader)

	application := Application{}
	config := config{}
	if err := settingsLoader.Load(&config); err != nil {
		logger.GetInstance().Error(
			"can't load application config: %s",
			zap.Error(err),
		)
		return nil
	}
	creatureTypesDB := db.New(&config.CreatureTypes.StorageType)
	if creatureTypesDB == nil {
		logger.GetInstance().Error("Can't create creature types mngr")
		return nil
	}
	application.ctx, application.ctxCancel = context.WithCancel(context.Background())
	application.routers = routers.New(application.ctx, settingsLoader, creatureTypesDB)
	if len(application.routers) < 1 {
		panic("no routers created - check env")
	}

	return &application
}

func (application *Application) Run() {
	defer application.ctxCancel()
	application.forEach(func(router irouter.Router) { go router.Run() })
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.forEach(func(router irouter.Router) { go router.Stop() })
}

func (application *Application) forEach(processer RouterProcess) {
	for _, router := range application.routers {
		if router == nil {
			panic("some of application routers is nil")
		}
		processer(router)
	}
}
