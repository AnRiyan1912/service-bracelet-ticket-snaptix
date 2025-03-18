package domain

type BraceletQrCodeData struct {
	NoTicket string                                `json:"no_ticket"`
	EventID  string                                `json:"event_id"`
	Sessions []BraceletSession                     `json:"sessions"`
	Category FindEventBraceletCategoryWithCategory `json:"category"`
}

type BraceletSession struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type BraceletQrCodeService interface {
	GenerateBraceletQrCode(braceletQrCodeData []BraceletQrCodeData, folderName string) (string, error)
}
