package publisher

import (
	"context"
	"fmt"
	"github.com/James226/kafka-accumulator/producer/app"
)

type WriterPublisher struct{}

func NewWriterPublisher() *WriterPublisher {
	return &WriterPublisher{}
}

func (WriterPublisher) Publish(ctx context.Context, msg app.DeltaMessage) {
	fmt.Printf("Message Received: %v\n", msg)
}

