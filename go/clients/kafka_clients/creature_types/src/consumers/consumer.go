package kafka_client_test_consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	creature_types_data "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/creature_types/data"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type ReaderConfig struct {
	Kafka struct {
		External struct {
			Servers []string
		}
		SendCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeConsumer struct {
	config ReaderConfig
	reader *kafka.Reader
}

func NewReader(settingsLoader isettings.SettingsLoader) *GetCreatureTypeConsumer {
	reader := GetCreatureTypeConsumer{}
	if err := settingsLoader.Load(&reader.config); err != nil {
		panic(fmt.Errorf("can't load handler config: %s", err))
	}
	reader.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:     reader.config.Kafka.External.Servers,
		Topic:       reader.config.Kafka.SendCreatureTypes.Topic,
		Partition:   0,
		StartOffset: 0,
	})
	logger.GetInstance().Info(
		"reader created",
		zap.String("reader", fmt.Sprintf("%v", reader)),
	)
	return &reader
}

func (consumer *GetCreatureTypeConsumer) Run(ctx context.Context, finishChan chan<- struct{}) {
	for {
		select {
		case <-ctx.Done():
			logger.GetInstance().Info(
				"stop consumer by context",
			)
			return
		default:
			msg, err := consumer.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
			}
			logger.GetInstance().Info(
				"get msg",
				zap.String("msg", fmt.Sprintf("%v", msg)),
				zap.String("msg.value", fmt.Sprintf("%v", msg.Value)),
				zap.String("msg.topic", fmt.Sprintf("%v", msg.Topic)),
			)
			var types creature_types_data.Types
			if err := json.Unmarshal(msg.Value, &types); err != nil {
				panic(fmt.Sprintf("can't read message: %s", err))
			}
			logger.GetInstance().Info(
				"msg processed",
				zap.String("types", fmt.Sprintf("%v", types)),
			)
			finishChan <- struct{}{}
			return
		}
	}
}
