package controllers

import (
	data_handlers "creature_types_srv/src/data_handlers/handlers"
	idata_handlers "creature_types_srv/src/data_handlers/interface"
	logger "creature_types_srv/src/logger/zap"
	"creature_types_srv/src/settings/implementations"

	"go.uber.org/zap"
)

type Controller struct {
	handler idata_handlers.IDataHandler
}

func (controller *Controller) Run() {
	settings := implementations.GetInstance().GetSettings()
	logger.GetInstance().Info(settings.PrettyJson())
	controller.handler = &data_handlers.DataHandler{}
	var err error
	if err = controller.handler.Init(); err != nil {
		logger.GetInstance().Error(
			"can't init data handler",
			zap.String("error", err.Error()))
		return
	}
	values := controller.handler.GetCreatureTypes()
	logger.GetInstance().Info(
		"Presented types",
		zap.String("types", values.String()),
	)
}
