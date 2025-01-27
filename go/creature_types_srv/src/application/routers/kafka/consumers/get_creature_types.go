package consumers

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	ikafka "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/interfaces"
	"github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/application/routers/kafka/utils"
	logger "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/logger/zap"
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
	consumer   *kafka.Consumer
	msgHandler ikafka.MsgHandler
	config     config
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewGetCreatureType(ctx context.Context, settingsLoader settings.SettingsLoader, handler ikafka.MsgHandler) (*GetCreatureType, error) {
	if handler == nil {
		return nil, fmt.Errorf("msg handler is nil")
	} else if settingsLoader == nil {
		return nil, fmt.Errorf("bad settings loader")
	}

	var getter GetCreatureType
	getter.msgHandler = handler
	if err := settingsLoader.Load(&getter.config); err != nil {
		return nil, fmt.Errorf("Can't load config: %w", err)
	}
	if getter.config.Kafka.GetCreatureTypes.IsGroupUnique {
		getter.config.Kafka.GetCreatureTypes.GroupID = uuid.New().String()
	} else if getter.config.Kafka.GetCreatureTypes.GroupID == "" {
		return nil, fmt.Errorf("unique group setted - setup group name")
	}
	getter.ctx, getter.cancel = context.WithCancel(ctx)

	var err error
	getter.consumer, err = kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": getter.config.Kafka.Servers,
			"group.id":          getter.config.Kafka.GetCreatureTypes.GroupID,
		})
	if err != nil {
		return nil, fmt.Errorf("Can't create consumer: %w", err)
	}

	if err = utils.CreateTopicByConsumer(
		ctx,
		settingsLoader,
		getter.consumer,
		&getter.config.Kafka.GetCreatureTypes.Topic,
	); err != nil {
		return nil, fmt.Errorf("can't create topic '%s': %w", getter.config.Kafka.GetCreatureTypes.Topic, err)
	}

	err = getter.consumer.SubscribeTopics(
		[]string{
			getter.config.Kafka.GetCreatureTypes.Topic,
		},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Consumer can't subscribe to topics: %w", err)
	}
	logger.GetInstance().Info(
		"consumer created",
		zap.String("consumer", fmt.Sprintf("%v", getter)),
	)
	return &getter, nil
}

func (getter *GetCreatureType) Run() error {
	if err := getter.isValid(); err != nil {
		return fmt.Errorf("trying to run invalid route: %w", err)
	}
	logger.GetInstance().Info(
		"consumer started",
		zap.String("config", fmt.Sprintf("%v", getter.config)),
	)
	for {
		select {
		case <-getter.ctx.Done():
			logger.GetInstance().Info(
				"stop consumer by context",
			)
			return nil
		default:
			event := getter.consumer.Poll(100)
			if event == nil {
				continue
			}

			logger.GetInstance().Info(
				"get kafka event",
				zap.String("event", fmt.Sprintf("%v", event)),
			)

			getter.handleEvent(event)
		}
	}
}

func (getter *GetCreatureType) Stop() error {
	getter.cancel()
	return getter.consumer.Close()
}

func (getter *GetCreatureType) handleEvent(event kafka.Event) {
	switch event := event.(type) {
	case *kafka.Message:
		logger.GetInstance().Info(
			"get message",
			zap.String("header", fmt.Sprintf("%v", event.Headers)),
		)
		go func() {
			err := getter.msgHandler.Handle(getter.ctx, event)
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
			zap.String("consumer", fmt.Sprintf("%v", *getter)),
		)
	}
}

func (consumer *GetCreatureType) isValid() error {
	if consumer.consumer == nil {
		return fmt.Errorf("bad reader")
	} else if consumer.msgHandler == nil {
		return fmt.Errorf("bad msg handler")
	}
	return nil
}
