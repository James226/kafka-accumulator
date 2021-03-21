package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	reader *kafka.Reader
}

func NewKafkaProducer(newReader func(kafka.ReaderConfig) *kafka.Reader) *KafkaProducer {
	l := log.New(os.Stdout, "kafka writer: ", 0)
	reader := newReader(kafka.ReaderConfig{
		Brokers:          []string{"kafka:9092"},
		GroupID:          "delta-accumulator",
		Topic:            "deltas",
		JoinGroupBackoff: time.Second,
		Logger:           l,
	})

	return &KafkaProducer{reader: reader}
}

func (k *KafkaProducer) ReadMessage(ctx context.Context) (bool, DeltaMessage) {
	select {
	case <-ctx.Done():
		return false, DeltaMessage{}
	default:
		m, err := k.reader.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Error reading from kafka: %v\n", err)
			return false, DeltaMessage{}
		}

		message := DeltaMessage{}
		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			fmt.Printf("Error reading from kafka: %v (%s)\n", err, string(m.Value))
			return false, DeltaMessage{}
		}

		fmt.Printf("Received message %v\n", message)

		return true, message
	}
}
