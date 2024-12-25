package consumers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	kafkaInterfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/kafka/interfaces"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type Config struct {
	Kafka struct {
		Servers          []string
		GetCreatureTypes struct {
			Topic         string
			IsGroupUnique bool
			GroupID       string `envconfig:"optional"`
		}
	}
}

type GetCreatureType struct {
	reader     *kafka.Reader
	msgHandler kafkaInterfaces.MsgHandler
	Config     Config
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewGetCreatureType(ctx context.Context, settings_loader settings.SettingsLoader) *GetCreatureType {
	consumer := GetCreatureType{}
	if err := settings_loader.Load(&consumer.Config); err != nil {
		panic("Bad settings loader")
	}
	if consumer.Config.Kafka.GetCreatureTypes.IsGroupUnique {
		consumer.Config.Kafka.GetCreatureTypes.GroupID = uuid.New().String()
	} else if consumer.Config.Kafka.GetCreatureTypes.GroupID == "" {
		panic("unique group setted - setup group name")
	}
	consumer.ctx, consumer.cancel = context.WithCancel(ctx)
	consumer.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: consumer.Config.Kafka.Servers,
		Topic:   consumer.Config.Kafka.GetCreatureTypes.Topic,
		GroupID: consumer.Config.Kafka.GetCreatureTypes.GroupID,
	})
	return &consumer
}

func (consumer *GetCreatureType) Run() {
	if err := consumer.isValid(); err != nil {
		panic(fmt.Errorf("trying to run invalid route: %s", err))
	}
	for {
		select {
		case <-consumer.ctx.Done():
			logger.GetInstance().Info(
				"stop consumer by context",
			)
			return
		default:
			msg, err := consumer.reader.ReadMessage(consumer.ctx)
			if err != nil {
				if consumer.ctx.Err() != nil {
					return
				}
			}
			go consumer.msgHandler.Handle(consumer.ctx, &msg)
		}
	}
}

func (consumer *GetCreatureType) Stop() {
	defer consumer.reader.Close()
	consumer.cancel()
}

func (consumer *GetCreatureType) SetupMsgHandler(handler kafkaInterfaces.MsgHandler) error {
	if handler == nil {
		return fmt.Errorf("msg handler is nil")
	}
	consumer.msgHandler = handler
	return nil
}

func (consumer *GetCreatureType) isValid() error {
	if consumer.reader == nil {
		return fmt.Errorf("bad reader")
	} else if consumer.msgHandler == nil {
		return fmt.Errorf("bad msg handler")
	}
	return nil
}
