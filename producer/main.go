package main

import (
	"context"
	"fmt"
	"github.com/James226/kafka-accumulator/producer/app"
	"github.com/James226/kafka-accumulator/producer/provider"
	"github.com/James226/kafka-accumulator/producer/publisher"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	prov := provider.NewRandomProvider()
	pub := publisher.NewKafkaPublisher(ctx)

	a := app.NewApp(prov, pub)

	a.Start(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	cancel()

	a.Join()
	fmt.Println("Shutdown complete")
}
