package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type TotalMessage struct {
	Id       string `json:"id"`
	Quantity int64  `json:"quantity"`
}

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(ctx context.Context) *KafkaPublisher {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:9092", "totals", 0)
	if err != nil {
		panic(err.Error())
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP("kafka:9092"),
		Topic:        "totals",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return &KafkaPublisher{writer: w}
}

func (p *KafkaPublisher) Publish(ctx context.Context, msg TotalMessage) error {
	bytes, _ := json.Marshal(msg)
	fmt.Printf("Writing: %s", string(bytes))
	message := kafka.Message{
		Key:   []byte(msg.Id),
		Value: bytes,
	}
	err := p.writer.WriteMessages(ctx, message)
	fmt.Printf("Written: %s", string(bytes))

	return err
}
