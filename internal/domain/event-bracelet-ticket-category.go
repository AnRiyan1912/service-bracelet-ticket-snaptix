package domain

type EventBraceletTicketCategory struct {
	ID             string `json:"id"`
	EventID        string `json:"event_id"`
	CategoryID     string `json:"category_id"`
	MaxUsePerEvent int    `json:"max_use_per_event"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type FindEventBraceletCategoryWithCategory struct {
	ID                string `json:"id"`
	CategoryName      string `json:"category_name"`
	IsUsedSeveralTime bool   `json:"is_used_several_time"`
	MaxUsePerEvent    int    `json:"max_use_per_event"`
}

type FindEventBraceletCategoryWithCategoryByEventID struct {
	ID           string `json:"id"`
	CategoryName string `json:"category_name"`
	MaxUse       int    `json:"max_use"`
}

type MysqlEventBraceletTicketCategoryRepository interface {
	InsertEventBraceletTicketCategory(eventBraceletTicketCategory EventBraceletTicketCategory) error
	FindByEventIDAndCategoryID(eventID string, categoryID string) (*FindEventBraceletCategoryWithCategory, error)
	FindMaxUsePerEventByEventIDAndCategoryID(ID string, eventID string) (int, error)
	FindAllByEventID(eventID string) ([]FindEventBraceletCategoryWithCategoryByEventID, error)
}
