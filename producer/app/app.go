package app

import (
	"context"
	"sync"
)

//go:generate mockgen -source=app.go -destination=../mocks/app.mocks.go
type Provider interface {
	Start(ctx context.Context) chan DeltaMessage
}

type Publisher interface {
	Publish(ctx context.Context, msg DeltaMessage)
}

type DeltaMessage struct {
	Id string
	Quantity int
}

type App struct {
	provider  Provider
	publisher Publisher
	finished  sync.WaitGroup
}

func NewApp(provider Provider, publisher Publisher) *App {
	return &App{
		provider: provider,
		publisher: publisher,
	}
}

func (a *App) Start(ctx context.Context) {
	providerChannel := a.provider.Start(ctx)
	a.finished = sync.WaitGroup{}
	a.finished.Add(1)

	go a.messageProcessor(ctx, providerChannel)
}

func (a *App) Join() {
	a.finished.Wait()
}

func (a *App) messageProcessor(ctx context.Context, providerChannel chan DeltaMessage) {
	for {
		select {
			case msg := <- providerChannel:
				a.publisher.Publish(ctx, msg)
			case <- ctx.Done():
				a.finished.Done()
				return
		}
	}
}