package controllers

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	data_handlers "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/handlers"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type RouterProcess func(router interfaces.Router)

type Controller struct {
	routers []interfaces.Router
}

func New() *Controller {
	loader := settings.New()
	if loader == nil {
		return nil
	}
	var settings_loader isettings.SettingsLoader = loader

	logger.InitLogger(settings_loader)

	handler := data_handlers.New(settings_loader)
	if handler == nil {
		logger.GetInstance().Error("can't create data handler")
		return nil
	}

	controller := Controller{}
	controller.routers = routers.New(settings_loader)
	if len(controller.routers) < 1 {
		panic("no routers created - check env")
	}
	controller.forEach(func(router interfaces.Router) {
		if err := router.Init(handler); err != nil {
			err_msg := fmt.Errorf("can't setup data handler: %s", err)
			panic(err_msg)
		}
	})

	return &controller
}

func (controller *Controller) Run() {
	controller.forEach(func(router interfaces.Router) { go router.Run() })
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	controller.forEach(func(router interfaces.Router) { go router.Stop() })
}

func (controller *Controller) forEach(processer RouterProcess) {
	for _, router := range controller.routers {
		if router == nil {
			panic("some of controller routers is nil")
		}
		processer(router)
	}
}
