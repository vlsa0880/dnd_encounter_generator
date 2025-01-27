package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers"
	irouter "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/interfaces"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/data_access/db"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func NewApplication() (*Application, error) {
	loader, err := settings.New()
	if err != nil {
		return nil, fmt.Errorf("can't construct settings loader: %w", err)
	}
	var settingsLoader isettings.SettingsLoader = loader

	logger.InitLogger(settingsLoader)

	application := Application{}
	config := config{}
	if err := settingsLoader.Load(&config); err != nil {
		return nil, fmt.Errorf("can't load application config: %w", err)
	}
	creatureTypesDB, err := db.New(&config.CreatureTypes.StorageType)
	if err != nil {
		return nil, fmt.Errorf("Can't create creature types mngr: %w", err)
	}
	application.ctx, application.ctxCancel = context.WithCancel(context.Background())

	application.routers, err = routers.New(application.ctx, settingsLoader, creatureTypesDB)
	if err != nil {
		return nil, fmt.Errorf("can't construct application: %w", err)
	}

	return &application, nil
}

func (application *Application) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	var wgFinishRun sync.WaitGroup
	go application.startWatchRun(stop, &wgFinishRun)
	<-stop
	application.waitRoutersStop()
	wgFinishRun.Wait()
}

func (application *Application) startWatchRun(stopChan chan<- os.Signal, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	var eg errgroup.Group
	application.forEach(func(router irouter.Router) { go eg.Go(router.Run) })
	if err := eg.Wait(); err != nil {
		logger.GetInstance().Error(
			"router run error",
			zap.Error(err),
		)
		stopChan <- syscall.SIGTERM
	}
}

func (application *Application) waitRoutersStop() {
	var eg errgroup.Group
	application.forEach(func(router irouter.Router) { go eg.Go(router.Stop) })
	if err := eg.Wait(); err != nil {
		logger.GetInstance().Error(
			"router stop error",
			zap.Error(err),
		)
	}
}

func (application *Application) forEach(processer RouterProcess) {
	for _, router := range application.routers {
		if router == nil {
			panic("some of application routers is nil")
		}
		processer(router)
	}
}
