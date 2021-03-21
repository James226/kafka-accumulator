package app_test

import (
	"context"
	a "github.com/James226/kafka-accumulator/accumulator/app"
	mock_app "github.com/James226/kafka-accumulator/accumulator/mocks"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestWhenTheAppIsStarted(t *testing.T) {
	done := sync.WaitGroup{}
	done.Add(1)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	producer := mock_app.NewMockProducer(ctl)
	processor := mock_app.NewMockDeltaProcessor(ctl)
	app := a.NewApp(producer, processor)

	ctx, cancel := context.WithCancel(context.Background())

	complete := sync.WaitGroup{}
	complete.Add(1)

	message := a.DeltaMessage{}
	producer.EXPECT().ReadMessage(ctx).Return(true, message)
	processor.EXPECT().Process(ctx, message).Do(func(context.Context, a.DeltaMessage) { cancel() })

	app.Run(ctx)
}