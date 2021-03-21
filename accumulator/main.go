package main

import (
	"context"
	"fmt"
	"github.com/James226/kafka-accumulator/accumulator/app"
	"github.com/segmentio/kafka-go"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	application := app.NewApp(
		app.NewKafkaProducer(kafka.NewReader),
		app.NewProcessor(app.NewKafkaPublisher(ctx)),
		)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	application.Run(ctx)

	fmt.Println("Shutdown complete")
}
