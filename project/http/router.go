package http

import (
	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(pub message.Publisher) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		Publisher: pub,
	}

	e.GET("/health", handler.Health)
	e.POST("/tickets-status", handler.PostTicketsStatus)

	return e
}
