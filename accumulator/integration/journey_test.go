// +build integration

package integration_test

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJourneyDeltaMessages(t *testing.T) {
	msg := "{\"id\": \"473285\", \"quantity\": 3}"

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "deltas", 0)
	assert.NoError(t, err)
	defer func() { _ = conn.Close() }()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "totals",
		GroupID: "TestJourneyDeltaMessages",
		StartOffset: kafka.LastOffset,
		JoinGroupBackoff: time.Second,
		MaxWait: time.Millisecond * 100,
	})
	_ = reader.SetOffset(kafka.LastOffset)

	defer func() { _ = reader.Close() }()

	message := kafka.Message{
		Key:   []byte("473285"),
		Value: []byte(msg),
	}

	_, err = conn.WriteMessages(message)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 25)
	defer cancel()
	receivedMessage, err := reader.ReadMessage(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "{\"id\":\"473285\",\"quantity\":3}", string(receivedMessage.Value))

	_, err = conn.WriteMessages(message)
	assert.NoError(t, err)

	receivedMessage, _ = reader.ReadMessage(ctx)

	assert.Equal(t, "{\"id\":\"473285\",\"quantity\":6}", string(receivedMessage.Value))
}

func TestJourneyTotalMessage(t *testing.T) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "deltas", 0)
	assert.NoError(t, err)
	defer func() { _ = conn.Close() }()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "totals",
		GroupID: "TestJourneyTotalMessage",
		StartOffset: kafka.LastOffset,
		JoinGroupBackoff: time.Second,
		MaxWait: time.Millisecond * 100,
	})
	_ = reader.SetOffset(kafka.LastOffset)

	defer func() { _ = reader.Close() }()
	message := kafka.Message{
		Key:   []byte("473287"),
		Value: []byte(("{\"id\": \"473287\", \"quantity\": 5}")),
	}

	_, err = conn.WriteMessages(message)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 25)
	defer cancel()
	receivedMessage, err := reader.ReadMessage(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "{\"id\":\"473287\",\"quantity\":5}", string(receivedMessage.Value))

	totalMessage := kafka.Message{
		Key:   []byte("473287"),
		Value: []byte("{\"id\": \"473287\", \"quantity\": 3, \"type\": \"total\"}"),
	}

	_, err = conn.WriteMessages(totalMessage)
	assert.NoError(t, err)

	receivedMessage, err = reader.ReadMessage(ctx)
	assert.NoError(t, err)

	assert.Equal(t, "{\"id\":\"473287\",\"quantity\":3}", string(receivedMessage.Value))
}