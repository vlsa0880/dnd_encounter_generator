package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type config struct {
	Kafka struct {
		Servers string
		Log     struct {
			Topic string
		}
	}
}

type LogSyncer struct {
	producer *kafka.Producer
	config   *config
}

func New(settingsLoader settings.SettingsLoader) (*LogSyncer, error) {
	if settingsLoader == nil {
		return nil, fmt.Errorf("Bad settings loader")
	}
	syncer := LogSyncer{}
	if err := settingsLoader.Load(&syncer.config); err != nil {
		return nil, fmt.Errorf("Can't load settings kafka log syncer config: %w", err)
	}
	var err error
	syncer.producer, err = kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": syncer.config.Kafka.Servers,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("Can't create producer: %w", err)
	}
	return &syncer, nil
}

func (syncer *LogSyncer) Write(msgData []byte) (int, error) {
	err := syncer.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &syncer.config.Kafka.Log.Topic,
				Partition: kafka.PartitionAny,
			},
			Value: msgData,
		},
		nil,
	)
	if err != nil {
		return 0, err
	}
	return len(msgData), nil
}
