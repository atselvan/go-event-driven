package entities

type TicketStatus string

var (
	TicketStatusConfirmed TicketStatus = "confirmed"
	TicketStatusCancelled TicketStatus = "canceled"
)

func (ts TicketStatus) String() string {
	return string(ts)
}
