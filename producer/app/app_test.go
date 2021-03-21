package app_test

import (
	"context"
	a "github.com/James226/kafka-accumulator/producer/app"
	mock_app "github.com/James226/kafka-accumulator/producer/mocks"
	"github.com/golang/mock/gomock"
	"sync"
	"testing"
)

func TestWhenTheAppIsStarted(t *testing.T) {
	providerChannel := make(chan a.DeltaMessage)
	done := sync.WaitGroup{}
	done.Add(1)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	provider := mock_app.NewMockProvider(ctl)
	publisher := mock_app.NewMockPublisher(ctl)
	app := a.NewApp(provider, publisher)

	ctx, cancel := context.WithCancel(context.Background())

	provider.EXPECT().Start(ctx).Return(providerChannel)

	app.Start(ctx)

	t.Run("a message is published", func(t *testing.T) {
		message := a.DeltaMessage{}
		publisher.EXPECT().Publish(ctx, message).Do(func(context.Context, a.DeltaMessage) { done.Done() })

		providerChannel <- message
		done.Wait()
	})

	cancel()
	app.Join()

}