package http

import (
	"net/http"
	"tickets/message"

	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		msg := message.Message{
			TicketID: ticket,
		}

		if err := h.broker.Send(msg); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
