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

type MysqlEventRepository interface {
	UpdateFolderNameTotalBraceletTicketByEventId(eventId string, folderName string, totalBraceletTicket int) error
	FindTotalAndTotalCheckInBraceletTicketByEventId(eventId string) (*Event, error)
}
