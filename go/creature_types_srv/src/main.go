package main

import (
	"os"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	utils "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/utils/common"
	"go.uber.org/zap"
)

func handleMainExit() {
	if runErr := recover(); runErr != nil {
		if errMsg, ok := runErr.(string); logger.GetInstance() != nil && ok {
			logger.GetInstance().Error(
				"panic occured",
				zap.String("errMsg", errMsg),
			)
		}
		if exit, ok := runErr.(utils.Exit); ok {
			os.Exit(exit.Code)
		}
		panic(runErr)
	}
}

func main() {
	defer handleMainExit()
	controller, err := application.NewApplication()
	if err != nil {
		logger.GetInstance().Error(
			"can't run application",
			zap.Error(err),
		)
	}
	controller.Run()
}
