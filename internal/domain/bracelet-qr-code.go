package domain

type BraceletQrCodeData struct {
	No       string `json:"no"`
	NoTicket string `json:"no_ticket"`
}

type BraceletQrCodeService interface {
	GenerateBraceletQrCode(braceletQrCodeData []BraceletQrCodeData) (string, error)
}
