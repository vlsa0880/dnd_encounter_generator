package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	controllers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/implementations"
	icontrollers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/controllers/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers"
	irouter "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/infrastructure/routers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/data_access/db"
	"go.uber.org/zap"
)

type RouterProcess func(router irouter.Router)

type Application struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	routers   []irouter.Router
	wgFinish  sync.WaitGroup
	stopCh    chan os.Signal
}

func NewApplication() (*Application, error) {
	loader, err := settings.New()
	if err != nil {
		return nil, fmt.Errorf("can't construct settings loader: %w", err)
	}
	var settingsLoader isettings.SettingsLoader = loader

	logger.InitLogger(settingsLoader)

	application := Application{}
	creatureTypesDB, err := db.New(settingsLoader)
	if err != nil {
		return nil, fmt.Errorf("Can't create creature types mngr: %w", err)
	}
	application.ctx, application.ctxCancel = context.WithCancel(context.Background())

	var controller icontrollers.CreatureTypes
	controller, err = controllers.NewGetCreatureTypes(creatureTypesDB)
	if err != nil {
		return nil, fmt.Errorf("Can't create creature types controller: %w", err)
	}

	application.routers, err = routers.New(application.ctx, settingsLoader, controller)
	if err != nil {
		return nil, fmt.Errorf("can't construct application: %w", err)
	}

	application.stopCh = make(chan os.Signal, 1)

	return &application, nil
}

func (application *Application) Run() {
	defer application.handlePanic()
	signal.Notify(application.stopCh, syscall.SIGTERM, syscall.SIGINT)
	application.forEachRouter(application.runRouter())
	<-application.stopCh
	application.forEachRouter(application.stopRouter())
	application.wgFinish.Wait()
	logger.GetInstance().Info("application run finished")
}

func (application *Application) runRouter() RouterProcess {
	return func(router irouter.Router) {
		go func() {
			defer application.handlePanic()
			defer application.wgFinish.Done()
			application.wgFinish.Add(1)
			if err := router.Run(); err != nil {
				logger.GetInstance().Error(
					"can't run router",
					zap.Error(err),
				)
				application.stopCh <- syscall.SIGTERM
			}
		}()
	}
}

func (application *Application) stopRouter() RouterProcess {
	return func(router irouter.Router) {
		go func() {
			defer application.handlePanic()
			if err := router.Stop(); err != nil {
				logger.GetInstance().Error(
					"can't stop router",
					zap.Error(err),
				)
			}
		}()
	}
}

func (application *Application) forEachRouter(processer RouterProcess) {
	for _, router := range application.routers {
		if router == nil {
			panic("some of application routers is nil")
		}
		processer(router)
	}
}

func (application *Application) handlePanic() {
	if runErr := recover(); runErr != nil {
		if errMsg, ok := runErr.(string); logger.GetInstance() != nil && ok {
			logger.GetInstance().Error(
				"panic occured - terminating",
				zap.String("errMsg", errMsg),
				zap.String("stacktrace", string(debug.Stack())),
			)
		}
		application.stopCh <- syscall.SIGTERM
	}
}
