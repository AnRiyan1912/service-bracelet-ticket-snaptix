package domain

type EncriptionData struct {
	NoTicketEncrypted string `json:"no_ticket_encrypted"`
	Nonce             string `json:"nonce"`
}
