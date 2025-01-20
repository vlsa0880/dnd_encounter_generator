package kafka_client_test_producer

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"

	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type WriterConfig struct {
	Kafka struct {
		External struct {
			Servers string
		}
		GetCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeProducer struct {
	config WriterConfig
	writer *kafka.Producer
}

func NewWriter(settingsLoader isettings.SettingsLoader) *GetCreatureTypeProducer {
	writer := GetCreatureTypeProducer{}
	var err error
	if err = settingsLoader.Load(&writer.config); err != nil {
		panic(fmt.Errorf("can't load handler config: %s", err))
	}
	logger.GetInstance().Info(
		"config successfully loaded",
		zap.String("config", fmt.Sprintf("%v", writer.config)),
	)
	writer.writer, err = kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": writer.config.Kafka.External.Servers,
		})
	if err != nil {
		panic(fmt.Sprintf("Can't create producer: %s", err))
	}
	logger.GetInstance().Info(
		"writer created",
		zap.String("writer", fmt.Sprintf("%v", writer)),
	)
	return &writer
}

func (producer *GetCreatureTypeProducer) Write(ctx context.Context) error {
	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &producer.config.Kafka.GetCreatureTypes.Topic,
			Partition: kafka.PartitionAny,
		},
		Headers: []kafka.Header{
			{Key: "uuid", Value: []byte(uuid.New().String())},
		},
	}
	logger.GetInstance().Info(
		"msg sended",
		zap.String("msg", fmt.Sprintf("%v", msg)),
	)
	produceConfirmation := make(chan kafka.Event)
	err := producer.writer.Produce(
		&msg,
		produceConfirmation,
	)
	logger.GetInstance().Info(
		"request sended",
	)
	select {
	case <-produceConfirmation:
		logger.GetInstance().Info(
			"msg produced successfully",
			zap.String("msg", fmt.Sprintf("%v", msg)),
		)
	case <-ctx.Done():
		logger.GetInstance().Error(
			"can't produce msg in time cause of timeout",
			zap.String("handler", fmt.Sprintf("%v", producer)),
		)
	}
	return err
}
