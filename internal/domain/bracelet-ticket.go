package domain

type BraceletTicket struct {
	ID                      string `json:"id"`
	EventID                 string `json:"event_id"`
	EventBraceletCategoryID string `json:"event_bracelet_category_id"`
	NoTicket                string `json:"no_ticket"`
	Status                  string `json:"status"`
	Sessions                string `json:"sessions"`
	SerialNumber            string `json:"serial_number"`
	NoTicketEncrypted       string `json:"no_ticket_encrypted"`
	DeviceID                string `json:"device_id"`
	DeviceName              string `json:"device_name"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
}

type BraceletSession struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
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

type GetBaceletTicketExelReq struct {
	EventID  string `json:"eventId" validate:"required"`
	FileName string `json:"fileName" validate:"required"`
}

type MysqlBraceletTicketRepository interface {
	InsertBraceletTicket(braceletTicket BraceletTicket) error
	FindByNoTicketEncrypted(noTicketEncrypted string) (*BraceletTicket, error)
	UpdateStatusDeviceIdAndNameById(ID string, deviceID string, deviceName string) error
	FindBySerialNumber(eventId string, serialNumber string) (*BraceletTicket, error)
	FindFirstWithLastSerialNumber(eventID string) (int, error)
}

type CheckInBraceletTicketOnlineRequest struct {
	EventID    string `json:"eventId" validate:"required"`
	QrData     string `json:"qrData" validate:"required"`
	DeviceID   string `json:"deviceId" validate:"required"`
	DeviceName string `json:"deviceName" validate:"required"`
}

type CheckInBraceletTicketOfflineRequest struct {
	Data []CheckInBraceletTicketOnlineRequest `json:"data" validate:"required"`
}

type CheckInBraceletTicketWithSerialNumberOnlineRequest struct {
	EventID      string `json:"eventId" validate:"required"`
	SerialNumber string `json:"serialNumber" validate:"required"`
	DeviceID     string `json:"deviceId" validate:"required"`
	DeviceName   string `json:"deviceName" validate:"required"`
}

type CheckInBraceletTicketWithSerialNumberOfflineRequest struct {
	Data []CheckInBraceletTicketWithSerialNumberOnlineRequest `json:"data" validate:"required"`
}

type BraceletTicketService interface {
	CheckInBraceletTicketOnline(eventId string, qrData string, deviceId string, deviceName string) (*ApiResponseWithaoutData, error)
	CheckInBraceletTicketOffline(data []CheckInBraceletTicketOnlineRequest) error
	GenerateBraceletQrCode(eventID string, braceletCategoryId string, total int, sessions []BraceletSession) error
	GetTotalBraceletAndTotalCheckInBraceletTicketByEventID(eventID string) (*GetTotalBraceletAndTotalCheckInBraceletTicketByEventIDRes, error)
	GetListFileNameExelBaceletTicketByEventID(eventID string) (*[]BraceletTicketExel, error)
	CheckInBraceletTicketOnlineManual(eventID string, serialNumber string, deviceID string, deviceName string) (*ApiResponseWithaoutData, error)
	CheckInBraceletTicketOfflineManual(data []CheckInBraceletTicketWithSerialNumberOnlineRequest) error
}
