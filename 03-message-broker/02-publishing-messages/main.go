package main

import (
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func main() {
	rdb := redis.NewClient(
		&redis.Options{
			Addr: os.Getenv("REDIS_ADDR"),
		},
	)

	logger := watermill.NewSlogLogger(nil)

	pub, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client: rdb,
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	msg := message.NewMessage(watermill.NewUUID(), []byte("50"))
	if err := pub.Publish("progress", msg); err != nil {
		panic(err)
	}

	msg = message.NewMessage(watermill.NewUUID(), []byte("100"))
	if err := pub.Publish("progress", msg); err != nil {
		panic(err)
	}
}
