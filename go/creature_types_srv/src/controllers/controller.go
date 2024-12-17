package controllers

import (
	data_handlers "creature_types_srv/src/data_handlers/handlers"
	"creature_types_srv/src/data_handlers/idata_handler"
	"creature_types_srv/src/http/routes/irouter"
	"creature_types_srv/src/http/routes/routers"
	logger "creature_types_srv/src/logger/zap"
	settings "creature_types_srv/src/settings/loader"

	"go.uber.org/zap"
)

type Controller struct {
	router irouter.IRouter
}

func New() *Controller {
	loader := settings.New()
	if loader == nil {
		return nil
	}
	var settings_loader settings.ISettingsLoader = loader

	logger.InitLogger(settings_loader)

	controller := Controller{}
	handler := data_handlers.New(settings_loader)
	if handler == nil {
		logger.GetInstance().Error("can't create data handler")
		return nil
	}
	var ihandler idata_handler.IDataHandler = handler
	router_mgr := routers.New(settings_loader)
	if router_mgr == nil {
		logger.GetInstance().Error("can't init router")
		return nil
	}
	controller.router = router_mgr

	var err error
	if err = controller.router.SetupDataHandler(ihandler); err != nil {
		logger.GetInstance().Error(
			"can't setup data handler",
			zap.String("err_msg", err.Error()),
		)
		return nil
	}
	return &controller
}

func (controller *Controller) Run() {
	controller.router.Run()
}
