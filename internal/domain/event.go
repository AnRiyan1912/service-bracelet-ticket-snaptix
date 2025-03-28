package domain

type Event struct {
	ID                         int    `json:"id"`
	EventName                  string `json:"event_name"`
	EventSessions              string `json:"event_sessions"`
	TotalBraceletTicket        int    `json:"total_bracelet_ticket"`
	TotalCheckInBraceletTicket int    `json:"total_check_in_bracelet_ticket"`
	CreatedAt                  string `json:"created_at"`
	UpdatedAt                  string `json:"updated_at"`
}

type EventTotalBraceletAndCheckInBraceletTicket struct {
	TotalBraceletTicket        int
	TotalCheckInBraceletTicket int
}

type MysqlEventRepository interface {
	UpdateTotalBraceletTicketByEventId(eventId string, totalBraceletTicket int) error
	UpdateTotalCheckInBraceletTicketByEventId(eventId string, totalCheckInBraceletTicket int) error
	FindTotalAndTotalCheckInBraceletTicketByEventId(eventId string) (*Event, error)
}

type RedisEventRepository interface {
	InsertTotalAndTotalCheckInBraceletTicketByEventId(eventId string, totalBraceletTicket int, totalCheckInBraceletTicket int) error
	FindTotalAndTotalCheckInBraceletTicketByEventId(eventId string) (*EventTotalBraceletAndCheckInBraceletTicket, error)
	UpdateTotalCheckInBraceletTicketByEventId(eventId string, totalCheckInBraceletTicket int) error
}
