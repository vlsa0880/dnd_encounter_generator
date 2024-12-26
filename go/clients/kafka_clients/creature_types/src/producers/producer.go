package kafka_client_test_producer

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type WriterConfig struct {
	Kafka struct {
		External struct {
			Servers []string
		}
		ProduceTimeout   time.Duration
		GetCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeProducer struct {
	config WriterConfig
	writer *kafka.Writer
}

func NewWriter(settingsLoader isettings.SettingsLoader) *GetCreatureTypeProducer {
	writer := GetCreatureTypeProducer{}
	if err := settingsLoader.Load(&writer.config); err != nil {
		panic(fmt.Errorf("can't load handler config: %s", err))
	}
	writer.writer = kafka.NewWriter(
		kafka.WriterConfig{
			Brokers: writer.config.Kafka.External.Servers,
			Topic:   writer.config.Kafka.GetCreatureTypes.Topic,
		})
	logger.GetInstance().Info(
		"writer created",
		zap.String("writer", fmt.Sprintf("%v", writer)),
	)
	return &writer
}

func (producer *GetCreatureTypeProducer) Write(ctx context.Context) error {
	msg := kafka.Message{
		Headers: []kafka.Header{
			{Key: "uuid", Value: []byte(uuid.New().String())},
		},
	}
	err := producer.writer.WriteMessages(
		ctx,
		msg,
	)
	logger.GetInstance().Info(
		"request sended",
	)
	return err
}
