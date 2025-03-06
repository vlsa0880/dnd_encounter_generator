package msghandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/utils"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
	"go.uber.org/zap"
)

type config struct {
	Kafka struct {
		Servers string
		Handle  struct {
			Timeout time.Duration
		}
		SendCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeMsgHandler struct {
	config        config
	creatureTypes dbinterfaces.CreatureTypesRepository
	producer      *kafka.Producer
}

func NewCreatureTypeHandler(settingsLoader settings.SettingsLoader, creatureTypes dbinterfaces.CreatureTypesRepository) (*GetCreatureTypeMsgHandler, error) {
	if creatureTypes == nil {
		return nil, fmt.Errorf("bad data handler")
	} else if settingsLoader == nil {
		return nil, fmt.Errorf("bad settings loader")
	}

	var handler GetCreatureTypeMsgHandler
	handler.creatureTypes = creatureTypes

	if err := settingsLoader.Load(&handler.config); err != nil {
		return nil, fmt.Errorf("can't load handler config: %w", err)
	}

	var err error
	handler.producer, err = kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": handler.config.Kafka.Servers,
		})
	if err != nil {
		return nil, fmt.Errorf("Can't create producer: %w", err)
	}

	if err = utils.CreateTopicByProducer(
		context.Background(),
		settingsLoader,
		handler.producer,
		&handler.config.Kafka.SendCreatureTypes.Topic,
	); err != nil {
		return nil, fmt.Errorf("can't create topic '%s': %w", handler.config.Kafka.SendCreatureTypes.Topic, err)
	}
	logger.GetInstance().Info("handler successfully created")
	return &handler, nil
}

func (handler *GetCreatureTypeMsgHandler) Handle(ctx context.Context, msg *kafka.Message) error {
	handleCtx, cancel := context.WithTimeout(ctx, handler.config.Kafka.Handle.Timeout)
	defer cancel()

	types, err := handler.creatureTypes.GetCreatureTypes(handleCtx)
	if err != nil {
		return fmt.Errorf("can't get creature types: %w", err)
	}
	logger.GetInstance().Debug(
		"types to send",
		zap.String("types", fmt.Sprintf("%v", types)),
	)

	typesData, err := json.Marshal(types)
	if err != nil {
		return fmt.Errorf("can't transform data for sending: %s", err)
	}
	logger.GetInstance().Debug(
		"types data to send",
		zap.String("types", fmt.Sprintf("%v", typesData)),
	)

	outMsg := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &handler.config.Kafka.SendCreatureTypes.Topic,
			Partition: kafka.PartitionAny,
		},
		Value:   typesData,
		Headers: msg.Headers,
	}
	logger.GetInstance().Info(
		"trying to produce message",
		zap.String("msg", fmt.Sprintf("%v", outMsg)),
	)

	produceConfirmation := make(chan kafka.Event)
	if err = handler.producer.Produce(&outMsg, produceConfirmation); err != nil {
		return fmt.Errorf(
			"error trying to produce msg '%v' to topic '%s': %w",
			&outMsg,
			handler.config.Kafka.SendCreatureTypes.Topic,
			err,
		)
	}

	select {
	case <-produceConfirmation:
		logger.GetInstance().Info(
			"msg produced successfully",
			zap.String("msg", fmt.Sprintf("%v", outMsg)),
		)
	case <-ctx.Done():
		logger.GetInstance().Error(
			"can't produce msg in time cause of timeout",
			zap.String("handler", fmt.Sprintf("%v", handler)),
		)
	}
	return err
}
