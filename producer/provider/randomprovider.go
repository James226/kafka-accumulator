package provider

import (
	"context"
	"github.com/James226/kafka-accumulator/producer/app"
	"math"
	"math/rand"
	"time"
)

type RandomProvider struct{}

func NewRandomProvider() *RandomProvider {
	return &RandomProvider{}
}

func (RandomProvider) Start(ctx context.Context) chan app.DeltaMessage {
	channel := make(chan app.DeltaMessage)

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		ids := []string {
			"328967",
			"198528",
			"029234",
			"183495",
			"549823",
		}

		for {
			select {
				case <- ctx.Done():
					return
				case <- ticker.C:
					random := float64(rand.Intn(8) - 4)
					channel <- app.DeltaMessage{Id: ids[rand.Intn(len(ids))], Quantity: int(math.Copysign(math.Abs(random) + 1, random))}
			}
		}
	}()

	return channel
}

