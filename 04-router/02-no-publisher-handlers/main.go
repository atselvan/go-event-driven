package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewSlogLogger(nil)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router := message.NewDefaultRouter(logger)
	router.AddConsumerHandler(
		"temperature_handler",
		"temperature-fahrenheit",
		sub,
		temperatureHandler,
	)

	err = router.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

func temperatureHandler(msg *message.Message) error {
	fmt.Printf("Temperature read: %s\n", string(msg.Payload))
	return nil
}
