package entities

type ConsumerGroup string

var (
	ConsumerGroupIssueReceipt         ConsumerGroup = "issue-receipt"
	ConsumerGroupAppendTicket         ConsumerGroup = "append-ticket"
	ConsumerGroupAppendTicketToRefund ConsumerGroup = "append-ticket-to-refund"
)

func (cg ConsumerGroup) String() string {
	return string(cg)
}
