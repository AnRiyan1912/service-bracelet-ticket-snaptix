package domain

type BraceletTicketExel struct {
	ID                  string `json:"id"`
	FileName            string `json:"file_name"`
	EventID             string `json:"event_id"`
	TotalBraceletTicket int    `json:"total_bracelet_ticket"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
}

type MysqlBraceletTicketExelRepository interface {
	InsertBraceletTicketExel(braceletTicketExel BraceletTicketExel) error
	FindByEventID(eventID string) ([]BraceletTicketExel, error)
	DeleteByEventID(eventID string) error
}

type BraceletTicketExelService interface {
	GenerateBraceletTicketExel(datas []BraceletQrCodeData, eventID string, fileName string) error
	GetBraceletTicketExelByEventID(eventID string) ([]BraceletTicketExel, error)
	DeleteBraceletTicketExelByEventID(eventID string) error
}
