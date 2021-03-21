package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"

	"github.com/James226/kafka-accumulator/producer/app"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(ctx context.Context) *KafkaPublisher {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", "deltas", 0)
	if err != nil {
		panic(err.Error())
	}

	w := &kafka.Writer{
		Addr:     kafka.TCP("kafka:9092"),
		Topic:    "deltas",
		Balancer: &kafka.LeastBytes{},
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return &KafkaPublisher{writer: w}
}

func (p *KafkaPublisher) Publish(ctx context.Context, msg app.DeltaMessage) {
	bytes, _ := json.Marshal(msg)
	fmt.Printf("Writing: %s", string(bytes))
	message := kafka.Message{
		Key:   []byte(msg.Id),
		Value: bytes,
	}
	err := p.writer.WriteMessages(ctx, message)

	if err != nil {
		panic(err)
	}
}
