package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(
		&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
		},
	)

	logger := watermill.NewSlogLogger(nil)

	sub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := sub.Subscribe(context.Background(), "progress")
	if err != nil {
		panic(err)
	}

	for msg := range msgs {
		fmt.Printf("Message ID: %s - %s\n", msg.UUID, msg.Payload)
		msg.Ack()
	}
}
