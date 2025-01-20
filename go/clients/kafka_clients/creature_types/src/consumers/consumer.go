package kafka_client_test_consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/domain/entities"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type ReaderConfig struct {
	Kafka struct {
		External struct {
			Servers string
		}
		SendCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeConsumer struct {
	config         ReaderConfig
	reader         *kafka.Consumer
	readyToProcess atomic.Bool
}

func NewReader(settingsLoader isettings.SettingsLoader) *GetCreatureTypeConsumer {
	reader := GetCreatureTypeConsumer{}
	reader.readyToProcess.Store(false)
	var err error
	if err = settingsLoader.Load(&reader.config); err != nil {
		panic(fmt.Errorf("can't load handler config: %s", err))
	}
	reader.reader, err = kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": reader.config.Kafka.External.Servers,
			"group.id":          "test_get_creature_type_group_" + uuid.NewString(),
		})
	if err != nil {
		panic(fmt.Sprintf("Can't create consumer: %s", err))
	}
	err = reader.reader.Subscribe(
		reader.config.Kafka.SendCreatureTypes.Topic,
		reader.rebalanceCallback,
	)
	if err != nil {
		panic(fmt.Sprintf("Consumer can't subscribe to topics: %s", err))
	}
	logger.GetInstance().Info(
		"consumer created",
	)
	return &reader
}

func (consumer *GetCreatureTypeConsumer) Run(ctx context.Context, finishChan chan<- struct{}) {
	defer consumer.reader.Close()
	logger.GetInstance().Info(
		"consumer started",
		zap.String("config", fmt.Sprintf("%v", consumer)),
	)
	for {
		select {
		case <-ctx.Done():
			logger.GetInstance().Info(
				"stop consumer by context",
			)
			return
		default:
			kafkaEvent := consumer.reader.Poll(1)
			if kafkaEvent == nil {
				continue
			}

			logger.GetInstance().Info(
				"get kafka event",
				zap.String("event", fmt.Sprintf("%v", kafkaEvent)),
			)

			switch event := kafkaEvent.(type) {
			case *kafka.Message:
				consumer.processMessage(finishChan, event)
				return
			case kafka.Error:
				logger.GetInstance().Warn(
					"get error msg from consumer",
					zap.String("err_msg", event.Error()),
				)
			default:
				continue
			}
		}
	}
}

func (consumer *GetCreatureTypeConsumer) WaitConsumerReady(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("Consumer didn't connect to topic in time")
		default:
			if consumer.readyToProcess.Load() {
				return nil
			}
			time.Sleep(time.Second)
		}
	}
}

func (consumer *GetCreatureTypeConsumer) processMessage(finishChan chan<- struct{}, msg *kafka.Message) {
	logger.GetInstance().Info(
		"get msg",
		zap.String("msg", fmt.Sprintf("%v", msg)),
		zap.String("msg.value", fmt.Sprintf("%v", msg.Value)),
		zap.String("msg.topic", fmt.Sprintf("%v", msg.TopicPartition.Topic)),
	)
	var types entities.CreatureTypes
	if err := json.Unmarshal(msg.Value, &types); err != nil {
		panic(fmt.Sprintf("can't read message: %s", err))
	}
	logger.GetInstance().Debug(
		"msg processed",
		zap.String("types", fmt.Sprintf("%v", types)),
	)
	finishChan <- struct{}{}
}

func (consumer *GetCreatureTypeConsumer) setConsumerReady(readiness bool) {
	consumer.readyToProcess.Store(readiness)
}

func (consumer *GetCreatureTypeConsumer) rebalanceCallback(eventConsumer *kafka.Consumer, incomeEvent kafka.Event) error {
	switch event := incomeEvent.(type) {
	case kafka.AssignedPartitions:
		fmt.Println("Partitions assigned:", event.Partitions)
		if err := eventConsumer.Assign(event.Partitions); err != nil {
			logger.GetInstance().Error(
				"can't assign consumer to new partitins",
				zap.Error(err),
			)
		} else {
			consumer.setConsumerReady(true)
		}
	case kafka.RevokedPartitions:
		fmt.Println("Partitions revoked:", event.Partitions)
		if err := eventConsumer.Assign(event.Partitions); err != nil {
			logger.GetInstance().Error(
				"can't unassign partitions from consumer",
				zap.Error(err),
			)
		} else {
			consumer.setConsumerReady(false)
		}
	}
	return nil
}
