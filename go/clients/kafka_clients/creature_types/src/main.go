package main

import (
	"context"
	"fmt"
	"time"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/env"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"

	consumers "github.com/vlsa0880/dnd_encounter_generator/go/clients/kafka_clients/creature_types/src/consumers"
	producers "github.com/vlsa0880/dnd_encounter_generator/go/clients/kafka_clients/creature_types/src/producers"
)

func main() {
	loader := settings.New()
	if loader == nil {
		panic("can't create settings loader")
	}
	var settingsLoader isettings.SettingsLoader = loader

	logger.InitLogger(settingsLoader)

	finishChan := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	consumer := consumers.NewReader(settingsLoader)
	go consumer.Run(ctx, finishChan)

	producer := producers.NewWriter(settingsLoader)
	if err := producer.Write(ctx); err != nil {
		panic(fmt.Sprintf("can't write data: %s", err))
	}
	logger.GetInstance().Debug(
		"waiting for response message...",
	)
	select {
	case <-ctx.Done():
		logger.GetInstance().Error(
			"response message wasn't proccessed in time",
		)
	case <-finishChan:
		logger.GetInstance().Debug(
			"response message was successfully processed",
		)
	}
}
