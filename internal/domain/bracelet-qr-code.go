package domain

type BraceletQrCodeData struct {
	NoTicket string                           `json:"no_ticket"`
	EventID  string                           `json:"event_id"`
	Sessions []BraceletSession                `json:"sessions"`
	Category BraceletTicketCategoryQrCodeData `json:"category"`
}

type BraceletTicketCategoryQrCodeData struct {
	Name              string `json:"name"`
	IsUsedSeveralTime bool   `json:"is_used_several_time"`
	MaxUse            int    `json:"max_use"`
}

type BraceletSession struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type BraceletQrCodeService interface {
	GenerateBraceletQrCode(braceletQrCodeData []BraceletQrCodeData, folderName string) (string, error)
}
