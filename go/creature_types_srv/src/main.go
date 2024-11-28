package main

import (
	"dnd_encounter_builder/src/logger"
	"dnd_encounter_builder/src/storage/settings"
)

func main() {
	settings := settings.GetInstance().GetSettings()
	logger.GetInstance().Info(settings.GetSlogGroup().String())
}
