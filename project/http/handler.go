package http

import (
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Publisher message.Publisher
}

func (h Handler) Health(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
