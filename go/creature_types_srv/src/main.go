package main

import (
	"os"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/controllers"
)

type Exit struct{ Code int }

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e)
	}
}

func main() {
	defer handleExit()
	controller := controllers.New()
	if controller == nil {
		panic("Can't create controller")
	}
	controller.Run()
}
