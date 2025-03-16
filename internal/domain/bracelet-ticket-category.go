package domain

type BraceletTicketCategory struct {
	ID                string `json:"id"`
	Name              string `json:"name"`
	IsUsedSeveralTime bool   `json:"is_used_several_time"`
	MaxUse            int    `json:"max_use"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type MysqlBraceletTicketCategoryRepository interface {
	InsertBraceletTicketCategory(braceletTicketCategory BraceletTicketCategory) error
	FindByBraceletTicketCategoryID(id string) (*BraceletTicketCategory, error)
	UpdateBraceletTicketCategory(braceletTicketCategory BraceletTicketCategory) error
}
