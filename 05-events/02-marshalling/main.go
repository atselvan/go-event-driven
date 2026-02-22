package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type PaymentCompleted struct {
	PaymentID   string `json:"payment_id"`
	OrderID     string `json:"order_id"`
	CompletedAt string `json:"completed_at"`
}

type OrderConfirmed struct {
	OrderId     string `json:"order_id"`
	ConfirmedAt string `json:"confirmed_at"`
}

func main() {
	logger := watermill.NewSlogLogger(nil)

	router := message.NewDefaultRouter(logger)

	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	sub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddHandler(
		"payment-completed-handler",
		"payment-completed",
		sub,
		"order-confirmed",
		pub,
		func(msg *message.Message) ([]*message.Message, error) {
			event := new(PaymentCompleted)
			err := json.Unmarshal(msg.Payload, event)
			if err != nil {
				return nil, err
			}

			newEvent := OrderConfirmed{
				OrderId:     event.OrderID,
				ConfirmedAt: event.CompletedAt,
			}

			payload, err := json.Marshal(&newEvent)
			if err != nil {
				return nil, err
			}

			outMsg := message.NewMessage(watermill.NewUUID(), payload)
			return []*message.Message{outMsg}, nil
		},
	)

	err = router.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
