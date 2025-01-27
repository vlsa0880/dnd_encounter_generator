package main

import (
	"os"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	"go.uber.org/zap"
)

type Exit struct{ Code int }

func handleExit() {
	if e := recover(); e != nil {
		if errMsg, ok := e.(string); logger.GetInstance() != nil && ok {
			logger.GetInstance().Error(
				"panic occured",
				zap.String("errMsg", errMsg),
			)
		}
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e)
	}
}

func main() {
	defer handleExit()
	controller, err := application.NewApplication()
	if err != nil {
		logger.GetInstance().Error(
			"can't run application",
			zap.Error(err),
		)
	}
	controller.Run()
}
