package controllers

import (
	data_handlers "creature_types_srv/src/data_handlers/handlers"
	logger "creature_types_srv/src/logger/zap"
	"creature_types_srv/src/routes/interfaces"
	"creature_types_srv/src/routes/routers"
	settings "creature_types_srv/src/settings/loader"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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
	var settings_loader settings.ISettingsLoader = loader

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
		if err := router.SetupDataHandler(handler); err != nil {
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
