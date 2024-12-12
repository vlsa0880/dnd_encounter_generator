package main

import (
	"creature_types_srv/src/controllers"
	"os"
)

type Exit struct{ Code int }

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e)
	}
}

func main() {
	defer handleExit()
	controller := controllers.Controller{}
	controller.Run()
}
