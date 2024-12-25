package msghandlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	dataHandlersInterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/data_handlers/interfaces"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type Config struct {
	Kafka struct {
		Servers           []string
		ProduceTimeout    time.Duration
		SendCreatureTypes struct {
			Topic string
		}
	}
}

type GetCreatureTypeMsgHandler struct {
	config      Config
	dataHandler dataHandlersInterfaces.DataHandler
	writer      *kafka.Writer
}

func NewCreatureTypeHandler(settingsLoader settings.SettingsLoader, dataHandler dataHandlersInterfaces.DataHandler) *GetCreatureTypeMsgHandler {
	if dataHandler == nil {
		panic("bad data handler")
	} else if settingsLoader == nil {
		panic("bad settings loader")
	}

	handler := &GetCreatureTypeMsgHandler{
		dataHandler: dataHandler,
	}

	if err := settingsLoader.Load(&handler.config); err != nil {
		panic(fmt.Errorf("can't load handler config: %s", err))
	}

	handler.writer = kafka.NewWriter(
		kafka.WriterConfig{
			Brokers: handler.config.Kafka.Servers,
			Topic:   handler.config.Kafka.SendCreatureTypes.Topic,
		})

	return handler
}

func (handler *GetCreatureTypeMsgHandler) Handle(ctx context.Context, msg *kafka.Message) error {
	handleCtx, cancel := context.WithTimeout(ctx, handler.config.Kafka.ProduceTimeout)
	defer cancel()
	types := handler.dataHandler.GetCreatureTypes(handleCtx)
	typesData, err := json.Marshal(types)
	if err != nil {
		return fmt.Errorf("can't transform data for sending: %s", err)
	}
	outMsg := kafka.Message{
		Value:   typesData,
		Headers: msg.Headers,
	}
	err = handler.writer.WriteMessages(
		handleCtx,
		outMsg,
	)
	return err
}
