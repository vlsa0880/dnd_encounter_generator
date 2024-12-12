package controllers

import (
	data_handlers "creature_types_srv/src/data_handlers/handlers"
	"creature_types_srv/src/http/routes/irouter"
	"creature_types_srv/src/http/routes/routers"
	logger "creature_types_srv/src/logger/zap"
	"creature_types_srv/src/settings/implementations"

	"go.uber.org/zap"
)

type Controller struct {
	router irouter.IRouter
}

func (controller *Controller) Run() {
	settings := implementations.GetInstance().GetSettings()
	logger.GetInstance().Info(settings.PrettyJson())
	handler := &data_handlers.DataHandler{}
	var err error
	if err = handler.Init(); err != nil {
		logger.GetInstance().Error(
			"can't init data handler",
			zap.String("error", err.Error()))
		return
	}
	router_mgr := routers.GinManager{}
	controller.router = &router_mgr
	if err = controller.router.Init(); err != nil {
		logger.GetInstance().Error(
			"can't init router",
			zap.String("err_msg", err.Error()),
		)
		return
	}
	if err = controller.router.SetupDataHandler(handler); err != nil {
		logger.GetInstance().Error(
			"can't setup data handler",
			zap.String("err_msg", err.Error()),
		)
		return
	}
	controller.router.Run()
}
