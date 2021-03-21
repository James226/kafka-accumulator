package app

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Publisher interface {
	Publish(ctx context.Context, msg TotalMessage) error
}

type Processor struct {
	publisher Publisher
	cache     map[string]int
	rdb       *redis.Client
}

func NewProcessor(publisher Publisher) *Processor {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	return &Processor{publisher: publisher, cache: make(map[string]int), rdb: rdb}
}

func (p Processor) Process(ctx context.Context, message DeltaMessage) {
	pipe := p.rdb.TxPipeline()
	var i int64
	str, err := p.rdb.Get(ctx, message.Id).Result()

	if err == redis.Nil {
	} else if err != nil {
		panic(err)
	} else {
		i, err = strconv.ParseInt(str, 10, 64)
	}

	if message.Type == "total" {
		i = message.Quantity
		_, err = pipe.Set(ctx, message.Id, message.Quantity, 0).Result()
	} else {
		i += message.Quantity
		_, err = pipe.IncrBy(ctx, message.Id, message.Quantity).Result()
	}

	totalMessage := TotalMessage{
		Id:       message.Id,
		Quantity: i,
	}

	p.publisher.Publish(ctx, totalMessage)
	_, err = pipe.Exec(ctx)
}
