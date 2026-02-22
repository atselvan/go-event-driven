package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tickets/entities"

	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

type TicketsStatusRequest struct {
	Tickets []TicketStatusRequest `json:"tickets"`
}

type TicketStatusRequest struct {
	TicketID      string         `json:"ticket_id"`
	Status        string         `json:"status"`
	Price         entities.Money `json:"price"`
	CustomerEmail string         `json:"customer_email"`
}

func (h Handler) PostTicketsStatus(c echo.Context) error {
	var request TicketsStatusRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		switch ticket.Status {
		case entities.TicketStatusConfirmed.String():
			event := entities.TicketBookingConfirmed{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := watermillMessage.NewMessage(event.Header.ID, payload)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

			if err := h.Publisher.Publish(entities.EventTicketBookingConfirmed.String(), msg); err != nil {
				return err
			}
		case entities.TicketStatusCancelled.String():
			event := entities.TicketBookingCanceled{
				Header:        entities.NewMessageHeader(),
				TicketID:      ticket.TicketID,
				CustomerEmail: ticket.CustomerEmail,
				Price:         ticket.Price,
			}

			payload, err := json.Marshal(event)
			if err != nil {
				return err
			}

			msg := watermillMessage.NewMessage(event.Header.ID, payload)
			msg.Metadata.Set("correlation_id", c.Request().Header.Get("Correlation-ID"))

			if err := h.Publisher.Publish(entities.EventTicketBookingCanceled.String(), msg); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown ticket status: %s", ticket.Status)
		}
	}

	return c.NoContent(http.StatusOK)
}
