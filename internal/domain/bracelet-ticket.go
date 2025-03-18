package domain

type BraceletTicket struct {
	ID                      string `json:"id"`
	EventBraceletCategoryID string `json:"event_bracelet_category_id"`
	NoTicket                string `json:"no_ticket"`
	Status                  string `json:"status"`
	Sessions                string `json:"sessions"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
}

type GenerateBraceletTicketReq struct {
	EventID            string            `json:"eventId" validate:"required"`
	BraceletCategoryID string            `json:"braceletCategoryId" validate:"required"`
	Total              int               `json:"total" validate:"required"`
	Sessions           []BraceletSession `json:"sessions" validate:"required"`
}

type GetTotalBraceletAndTotalCheckInBraceletTicketByEventIDRes struct {
	TotalBraceletTicket int `json:"totalBraceletTicket"`
	TotalCheckIn        int `json:"totalCheckIn"`
}

type MysqlBraceletTicketRepository interface {
	InsertBraceletTicket(braceletTicket BraceletTicket) error
	FindByNoTicket(noTicket string) (*BraceletTicket, error)
	UpdateStatusById(ID string) error
}

type CheckInBraceletTicketRequest struct {
	EventID string `json:"eventId" validate:"required"`
	QrData  string `json:"qrData" validate:"required"`
}

type BraceletTicketService interface {
	InsertBraceletTicket(braceletTicket BraceletTicket) error
	FindByBraceletTicketID(id int) (*BraceletTicket, error)
	UpdateBraceletTicket(braceletTicket BraceletTicket) error
	CheckInBraceletTicket(eventId string, qrData string) (*ApiResponseWithaoutData, error)
	GenerateBraceletQrCode(eventID string, braceletCategoryId string, total int, sessions []BraceletSession) error
	GetTotalBraceletAndTotalCheckInBraceletTicketByEventID(eventID string) (*GetTotalBraceletAndTotalCheckInBraceletTicketByEventIDRes, error)
}
