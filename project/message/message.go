package message

import (
	"context"
	"fmt"

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

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketId string) error
}

func NewClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func NewPublisher(rdb *redis.Client, logger watermill.LoggerAdapter) (message.Publisher, error) {
	pub, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client: rdb,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}
	return pub, err
}

func NewRouter(
	rdb *redis.Client,
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
	logger watermill.LoggerAdapter,
) (*message.Router, error) {
	sub, err := redisstream.NewSubscriber(
		redisstream.SubscriberConfig{
			Client: rdb,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	router := message.NewDefaultRouter(logger)

	router.AddConsumerHandler(
		"issue-receipt-handler",
		TopicIssueReceipt.String(),
		sub,
		issueReceiptHandler(receiptsService),
	)

	router.AddConsumerHandler(
		"append-to-tracker-handler",
		TopicAppendToTracker.String(),
		sub,
		appendToTrackerHandler(spreadsheetsAPI),
	)

	return router, nil
}

func issueReceiptHandler(srv ReceiptsService) message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		ticketID := string(msg.Payload)
		fmt.Printf("Processing message: topic: %s msg-payload: %v\n", TopicIssueReceipt, ticketID)
		if err := srv.IssueReceipt(msg.Context(), ticketID); err != nil {
			return err
		}
		return nil
	}
}

func appendToTrackerHandler(srv SpreadsheetsAPI) message.NoPublishHandlerFunc {
	return func(msg *message.Message) error {
		ticketID := string(msg.Payload)
		fmt.Printf("Processing message: topic: %s msg-payload: %v\n", TopicAppendToTracker, ticketID)
		if err := srv.AppendRow(msg.Context(), "tickets-to-print", []string{ticketID}); err != nil {
			return err
		}
		return nil
	}
}
