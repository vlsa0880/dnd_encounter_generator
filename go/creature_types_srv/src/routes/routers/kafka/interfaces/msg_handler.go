package interfaces

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type MsgHandler interface {
	Handle(ctx context.Context, msg *kafka.Message)
}
