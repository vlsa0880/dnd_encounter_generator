package consumers

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
	kafka_interfaces "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/routes/routers/kafka/interfaces"
	settings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
	"go.uber.org/zap"
)

type config struct {
	Kafka struct {
		Servers          string
		GetCreatureTypes struct {
			Topic         string
			IsGroupUnique bool
			GroupID       string `envconfig:"optional"`
		}
	}
}

type GetCreatureType struct {
	reader     *kafka.Consumer
	msgHandler kafka_interfaces.MsgHandler
	config     config
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewGetCreatureType(ctx context.Context, settingsLoader settings.SettingsLoader, handler kafka_interfaces.MsgHandler) *GetCreatureType {
	if handler == nil {
		panic("msg handler is nil")
	} else if settingsLoader == nil {
		panic("bad settings loader")
	}
	consumer := GetCreatureType{}
	consumer.msgHandler = handler
	if err := settingsLoader.Load(&consumer.config); err != nil {
		panic(fmt.Sprintf("Can't load config: %s", err))
	}
	if consumer.config.Kafka.GetCreatureTypes.IsGroupUnique {
		consumer.config.Kafka.GetCreatureTypes.GroupID = uuid.New().String()
	} else if consumer.config.Kafka.GetCreatureTypes.GroupID == "" {
		panic("unique group setted - setup group name")
	}
	consumer.ctx, consumer.cancel = context.WithCancel(ctx)
	var err error
	consumer.reader, err = kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": consumer.config.Kafka.Servers,
			"group.id":          consumer.config.Kafka.GetCreatureTypes.GroupID,
		})
	if err != nil {
		panic(fmt.Sprintf("Can't create consumer: %s", err))
	}
	err = consumer.reader.SubscribeTopics(
		[]string{
			consumer.config.Kafka.GetCreatureTypes.Topic,
		},
		nil,
	)
	if err != nil {
		panic(fmt.Sprintf("Consumer can't subscribe to topics: %s", err))
	}
	logger.GetInstance().Info(
		"consumer created",
		zap.String("consumer", fmt.Sprintf("%v", consumer)),
	)
	return &consumer
}

func (consumer *GetCreatureType) Run() {
	if err := consumer.isValid(); err != nil {
		panic(fmt.Errorf("trying to run invalid route: %s", err))
	}
	logger.GetInstance().Info(
		"consumer started",
		zap.String("config", fmt.Sprintf("%v", consumer.config)),
	)
	for {
		select {
		case <-consumer.ctx.Done():
			logger.GetInstance().Info(
				"stop consumer by context",
			)
			return
		default:
			event := consumer.reader.Poll(100)
			if event == nil {
				continue
			}

			logger.GetInstance().Info(
				"get kafka event",
				zap.String("event", fmt.Sprintf("%v", event)),
			)

			consumer.handleEvent(event)
		}
	}
}

func (consumer *GetCreatureType) Stop() {
	defer consumer.reader.Close()
	consumer.cancel()
}

func (consumer *GetCreatureType) handleEvent(event kafka.Event) {
	switch event := event.(type) {
	case *kafka.Message:
		logger.GetInstance().Info(
			"get message",
			zap.String("header", fmt.Sprintf("%v", event.Headers)),
		)
		go func() {
			err := consumer.msgHandler.Handle(consumer.ctx, event)
			if err != nil {
				logger.GetInstance().Error(
					"can't handle msg",
					zap.Error(err),
				)
			} else {
				logger.GetInstance().Info(
					"successfully handled event",
					zap.String("event", fmt.Sprintf("%v", event)),
				)
			}
		}()
	case kafka.Error:
		logger.GetInstance().Warn(
			"get error msg from consumer",
			zap.String("err_msg", event.Error()),
		)
	}
}

func (consumer *GetCreatureType) isValid() error {
	if consumer.reader == nil {
		return fmt.Errorf("bad reader")
	} else if consumer.msgHandler == nil {
		return fmt.Errorf("bad msg handler")
	}
	return nil
}
