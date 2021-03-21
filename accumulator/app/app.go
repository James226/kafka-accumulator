package app

import (
	"context"
)

//go:generate mockgen -source=app.go -destination=../mocks/app.mocks.go
type Producer interface {
	ReadMessage(ctx context.Context) (bool, DeltaMessage)
}

type DeltaProcessor interface {
	Process(ctx context.Context, message DeltaMessage)
}

type DeltaMessage struct {
	Id string
	Quantity int64
	Type string
}

type App struct {
	producer  Producer
	processor DeltaProcessor
}

func NewApp(producer Producer, processor DeltaProcessor) *App {
	return &App{ producer: producer, processor: processor }
}

func (a *App) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		default:
			received, msg := a.producer.ReadMessage(ctx)
			if received {
				a.processor.Process(ctx, msg)
			}

		}
	}
}