package service

import (
	"context"
	"errors"
	stdHTTP "net/http"
	"tickets/message"

	"github.com/labstack/echo/v4"

	ticketsHttp "tickets/http"
)

type Service struct {
	echoRouter *echo.Echo
	broker     *message.Broker
}

func New(
	spreadsheetsAPI message.SpreadsheetsAPI,
	receiptsService message.ReceiptsService,
) (Service, error) {
	b, err := message.NewBroker(spreadsheetsAPI, receiptsService)
	if err != nil {
		return Service{}, err
	}

	return Service{
		echoRouter: ticketsHttp.NewHttpRouter(b),
		broker:     b,
	}, nil
}

func (s Service) Run(ctx context.Context) error {
	go func() {
		err := s.broker.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()

	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}

	return nil
}
