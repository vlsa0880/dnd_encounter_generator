package msghandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	dbinterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/use_cases/interfaces"
	"go.uber.org/zap"
)

type Config struct {
	Kafka struct {
		Servers           string
		ProduceTimeout    time.Duration
		SendCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeMsgHandler struct {
	config        Config
	creatureTypes dbinterfaces.CreatureTypes
	writer        *kafka.Producer
}

func NewCreatureTypeHandler(settingsLoader settings.SettingsLoader, creatureTypes dbinterfaces.CreatureTypes) *GetCreatureTypeMsgHandler {
	if creatureTypes == nil {
		panic("bad data handler")
	} else if settingsLoader == nil {
		panic("bad settings loader")
	}

	handler := &GetCreatureTypeMsgHandler{
		creatureTypes: creatureTypes,
	}

	if err := settingsLoader.Load(&handler.config); err != nil {
		panic(fmt.Errorf("can't load handler config: %s", err))
	}

	var err error
	handler.writer, err = kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": handler.config.Kafka.Servers,
		})
	if err != nil {
		panic(fmt.Sprintf("Can't create producer: %s", err))
	}
	logger.GetInstance().Info("handler successfully created")
	return handler
}

func (handler *GetCreatureTypeMsgHandler) Handle(ctx context.Context, msg *kafka.Message) error {
	handleCtx, cancel := context.WithTimeout(ctx, handler.config.Kafka.ProduceTimeout)
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
		"trying to send message",
		zap.String("msg", fmt.Sprintf("%v", outMsg)),
	)
	err = handler.writer.Produce(
		&outMsg,
		nil,
	)
	return err
}
