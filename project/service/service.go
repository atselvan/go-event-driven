package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	stdHTTP "net/http"
	"os"
	"os/signal"
	"tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	ticketsHttp "tickets/http"
)

type Service struct {
	echoRouter *echo.Echo
	msgRouter  *watermillMessage.Router
}

func New(
	rdb *redis.Client,
	spreadsheetsAPI message.SpreadsheetsAPI,
	receiptsService message.ReceiptsService,
) (*Service, error) {
	logger := watermill.NewSlogLogger(slog.Default())

	publisher, err := message.NewPublisher(rdb, logger)
	if err != nil {
		return nil, err
	}

	router, err := message.NewRouter(rdb, spreadsheetsAPI, receiptsService, logger)
	if err != nil {
		return nil, err
	}

	return &Service{
		echoRouter: ticketsHttp.NewHttpRouter(publisher),
		msgRouter:  router,
	}, nil
}

func (s Service) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errGrp, ctx := errgroup.WithContext(ctx)

	errGrp.Go(func() error {
		return s.msgRouter.Run(ctx)
	})

	errGrp.Go(func() error {
		<-s.msgRouter.Running()

		err := s.echoRouter.Start(":8080")
		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}
		return nil
	})

	errGrp.Go(func() error {
		<-ctx.Done()
		fmt.Println("Shutting down the application...")
		return s.echoRouter.Shutdown(ctx)
	})

	if err := errGrp.Wait(); err != nil {
		return err
	}

	return nil
}
