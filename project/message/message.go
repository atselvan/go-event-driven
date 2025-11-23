package message

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type Topic string

var (
	TopicIssueReceipt    Topic = "issue-receipt"
	TopicAppendToTracker Topic = "append-to-tracker"
)

func (t Topic) String() string {
	return string(t)
}

type Message struct {
	TicketID string
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketId string) error
}

type Broker struct {
	client *redis.Client
	pub    message.Publisher

	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
}

func NewBroker(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) (*Broker, error) {
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
		return nil, err
	}

	return &Broker{
		client:          rdb,
		pub:             pub,
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}, nil
}

func (b *Broker) Send(m Message) error {
	msg := message.NewMessage(watermill.NewUUID(), []byte(m.TicketID))

	if err := b.pub.Publish(TopicIssueReceipt.String(), msg); err != nil {
		return err
	}

	if err := b.pub.Publish(TopicAppendToTracker.String(), msg); err != nil {
		return err
	}

	return nil
}

func (b *Broker) Run(ctx context.Context) error {
	logger := watermill.NewSlogLogger(nil)

	sub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: b.client,
		},
		logger,
	)
	if err != nil {
		return err
	}

	router := message.NewDefaultRouter(logger)

	router.AddConsumerHandler(
		"issue-receipt-handler",
		TopicIssueReceipt.String(),
		sub,
		b.issueReceiptHandler(ctx),
	)

	router.AddConsumerHandler(
		"append-to-tracker-handler",
		TopicAppendToTracker.String(),
		sub,
		b.appendToTrackerHandler(ctx),
	)

	go func() {
		err := router.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (b *Broker) issueReceiptHandler(ctx context.Context) message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		ticketID := string(msg.Payload)
		fmt.Printf("Processing message: topic: %s msg-payload: %v\n", TopicIssueReceipt, ticketID)
		if err := b.receiptsService.IssueReceipt(ctx, ticketID); err != nil {
			return err
		}
		return nil
	}
}

func (b *Broker) appendToTrackerHandler(ctx context.Context) message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		ticketID := string(msg.Payload)
		fmt.Printf("Processing message: topic: %s msg-payload: %v\n", TopicAppendToTracker, ticketID)
		if err := b.spreadsheetsAPI.AppendRow(ctx, "tickets-to-print", []string{ticketID}); err != nil {
			return err
		}
		return nil
	}
}
