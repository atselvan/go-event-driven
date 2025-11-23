package http

import (
	"net/http"
	"tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
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
		msg := watermillMessage.NewMessage(watermill.NewUUID(), []byte(ticket))

		if err := h.Publisher.Publish(message.TopicIssueReceipt.String(), msg); err != nil {
			return err
		}

		if err := h.Publisher.Publish(message.TopicAppendToTracker.String(), msg); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
