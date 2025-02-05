package interfaces

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MsgHandler interface {
	Handle(ctx context.Context, msg *kafka.Message) error
}
