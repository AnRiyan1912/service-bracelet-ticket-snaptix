package domain

type Ticket struct {
	EventSessions string `json:"event_sessions"`
}

type MysqlTicketRepository interface {
	FindEventSessionsByTicketIdID(ticketID string) (*Ticket, error)
}
