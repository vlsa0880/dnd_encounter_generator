package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	isettings "github.com/vlsa0880/dnd_encounter_generator/go/creature_types_srv/src/settings/loader/interfaces"
)

type config struct {
	Kafka struct {
		Auto struct {
			TopicsCreation        bool
			TopicsCreationTimeout time.Duration
		}
	}
}

func CreateTopicByConsumer(ctx context.Context, settingsLoader isettings.SettingsLoader, consumer *kafka.Consumer, topics *string) error {
	if settingsLoader == nil {
		return fmt.Errorf("bad settings loader")
	} else if consumer == nil {
		return fmt.Errorf("bad consumer")
	}

	config := config{}
	if err := settingsLoader.Load(&config); err != nil {
		return fmt.Errorf("can't load settings: %w", err)
	}

	topicCreationCtx, cancel := context.WithTimeout(ctx, config.Kafka.Auto.TopicsCreationTimeout)
	defer cancel()
	if config.Kafka.Auto.TopicsCreation {
		adminClient, err := kafka.NewAdminClientFromConsumer(consumer)
		if err != nil {
			return fmt.Errorf("Can't create admin from consumer: %s", err)
		}
		if err = CreateTopic(topicCreationCtx, adminClient, topics); err != nil {
			return fmt.Errorf("Can't create topic '%v': %s", topics, err)
		}
	}
	return nil
}

func CreateTopicByProducer(ctx context.Context, settingsLoader isettings.SettingsLoader, producer *kafka.Producer, topics *string) error {
	if settingsLoader == nil {
		return fmt.Errorf("bad settings loader")
	} else if producer == nil {
		return fmt.Errorf("bad consumer")
	}

	config := config{}
	if err := settingsLoader.Load(&config); err != nil {
		return fmt.Errorf("can't load settings: %w", err)
	}

	topicCreationCtx, cancel := context.WithTimeout(ctx, config.Kafka.Auto.TopicsCreationTimeout)
	defer cancel()
	if config.Kafka.Auto.TopicsCreation {
		adminClient, err := kafka.NewAdminClientFromProducer(producer)
		if err != nil {
			return fmt.Errorf("Can't create admin from consumer: %s", err)
		}
		if err = CreateTopic(topicCreationCtx, adminClient, topics); err != nil {
			return fmt.Errorf("Can't create topic '%v': %s", topics, err)
		}
	}
	return nil
}

func CreateTopic(ctx context.Context, adminClient *kafka.AdminClient, topic *string) error {
	exists, err := checkIfTopicExists(ctx, adminClient, topic)
	if err != nil {
		return err
	} else if exists {
		fmt.Printf("topic '%s' already exists", *topic)
		return nil
	}
	return createTopic(ctx, adminClient, topic)
}

func checkIfTopicExists(ctx context.Context, adminClient *kafka.AdminClient, topic *string) (bool, error) {
	timeToDeadline, _ := ctx.Deadline()
	timeoutMs := int(time.Until(timeToDeadline).Milliseconds())
	metadata, err := adminClient.GetMetadata(topic, false, timeoutMs)
	if err != nil {
		return false, fmt.Errorf("failed to get metadata: %v", err)
	}

	for _, t := range metadata.Topics {
		if t.Topic == *topic {
			return true, nil
		}
	}
	return false, nil
}

func createTopic(ctx context.Context, adminClient *kafka.AdminClient, topic *string) error {
	topics := []kafka.TopicSpecification{
		{
			Topic: *topic,
		},
	}

	results, err := adminClient.CreateTopics(ctx, topics)
	if err != nil {
		return fmt.Errorf("error creating topic: %v", err)
	}

	for _, res := range results {
		if res.Error.Code() != kafka.ErrNoError {
			return fmt.Errorf("failed to create topic %s: %v", *topic, res.Error)
		}
	}
	fmt.Printf("topic '%s' successfully created", *topic)
	return nil
}
