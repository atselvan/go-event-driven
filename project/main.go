package main

import (
	"log/slog"
	"os"
	"tickets/message"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	rdb := message.NewClient(os.Getenv("REDIS_ADDR"))

	spreadsheetsAPI := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	srv, err := service.New(rdb, spreadsheetsAPI, receiptsService)
	if err != nil {
		panic(err)
	}

	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
