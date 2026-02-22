package events

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error
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
	issueReceiptSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: entities.ConsumerGroupIssueReceipt.String(),
	}, logger)
	if err != nil {
		panic(err)
	}

	appendToTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: entities.ConsumerGroupAppendTicket.String(),
	}, logger)
	if err != nil {
		panic(err)
	}

	appendToRefundTrackerSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: entities.ConsumerGroupAppendTicketToRefund.String(),
	}, logger)
	if err != nil {
		panic(err)
	}

	router := message.NewDefaultRouter(logger)
	useMiddleware(router, logger)

	handler := NewHandler(spreadsheetsAPI, receiptsService)

	router.AddConsumerHandler(
		"issue-receipt-handler",
		entities.EventTicketBookingConfirmed.String(),
		issueReceiptSub,
		handler.IssueReceipt(),
	)

	router.AddConsumerHandler(
		"append-to-tracker-handler",
		entities.EventTicketBookingConfirmed.String(),
		appendToTrackerSub,
		handler.AppendToTracker(),
	)

	router.AddConsumerHandler(
		"append-to-refund-tracker-handler",
		entities.EventTicketBookingCanceled.String(),
		appendToRefundTrackerSub,
		handler.AppendTicketsToRefundTracker(),
	)

	return router, nil
}
