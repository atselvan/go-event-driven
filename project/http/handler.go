package http

import (
	"tickets/message"
)

type Handler struct {
	broker *message.Broker
}
