package entities

import (
	"time"

	"github.com/ThreeDotsLabs/watermill"
)

type Event string

var (
	EventTicketBookingConfirmed Event = "TicketBookingConfirmed"
	EventTicketBookingCanceled  Event = "TicketBookingCanceled"
)

func (e Event) String() string {
	return string(e)
}

type MessageHeader struct {
	ID          string    `json:"id"`
	PublishedAt time.Time `json:"published_at"`
}

func NewMessageHeader() MessageHeader {
	return MessageHeader{
		ID:          watermill.NewUUID(),
		PublishedAt: time.Now(),
	}
}

type TicketBookingConfirmed struct {
	Header MessageHeader

	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}

type TicketBookingCanceled struct {
	Header MessageHeader

	TicketID      string `json:"ticket_id"`
	CustomerEmail string `json:"customer_email"`
	Price         Money  `json:"price"`
}
