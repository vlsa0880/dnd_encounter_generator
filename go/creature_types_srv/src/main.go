package main

import (
	"creature_types_srv/src/controllers"
	"os"
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
