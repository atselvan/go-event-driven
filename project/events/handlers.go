package events

import (
	"encoding/json"
	"log/slog"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Handler struct {
	spreadsheetsAPI SpreadsheetsAPI
	receiptsAPI     ReceiptsService
}

func NewHandler(spreadsheetsAPI SpreadsheetsAPI, receiptsAPI ReceiptsService) Handler {
	return Handler{
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsAPI:     receiptsAPI,
	}
}

func (h Handler) AppendToTracker() message.NoPublishHandlerFunc {
	slog.Info("Appending ticket to the tracker")
	return func(msg *message.Message) error {
		var payload entities.AppendToTrackerPayload
		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		if err := h.spreadsheetsAPI.AppendRow(
			msg.Context(),
			"tickets-to-print",
			[]string{
				payload.TicketID,
				payload.CustomerEmail,
				payload.Price.Amount,
				payload.Price.Currency,
			}); err != nil {
			return err
		}
		return nil
	}
}

func (h Handler) AppendTicketsToRefundTracker() message.NoPublishHandlerFunc {
	slog.Info("Appending ticket to the refund tracker")
	return func(msg *message.Message) error {
		var payload entities.AppendToTrackerPayload
		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		if err := h.spreadsheetsAPI.AppendRow(
			msg.Context(),
			"tickets-to-refund",
			[]string{
				payload.TicketID,
				payload.CustomerEmail,
				payload.Price.Amount,
				payload.Price.Currency,
			}); err != nil {
			return err
		}
		return nil
	}
}

func (h Handler) IssueReceipt() message.NoPublishHandlerFunc {
	slog.Info("Issue Receipt")
	return func(msg *message.Message) error {
		var payload entities.IssueReceiptPayload
		err := json.Unmarshal(msg.Payload, &payload)
		if err != nil {
			return err
		}

		if err := h.receiptsAPI.IssueReceipt(
			msg.Context(),
			entities.IssueReceiptRequest{
				TicketID: payload.TicketID,
				Price:    payload.Price,
			}); err != nil {
			return err
		}
		return nil
	}
}
