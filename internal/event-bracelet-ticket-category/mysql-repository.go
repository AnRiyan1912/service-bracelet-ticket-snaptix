package eventbraceletticketcategory

import (
	"gorm.io/gorm"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type mysqlEventBraceletTicketCategoryRepository struct {
	db *gorm.DB
}

func NewMysqlEventBraceletTicketCategoryRepository(db *gorm.DB) domain.MysqlEventBraceletTicketCategoryRepository {
	return &mysqlEventBraceletTicketCategoryRepository{db: db}
}

// FindByEventIDAndCategoryID implements domain.MysqlEventBraceletTicketCategoryRepository.
func (m *mysqlEventBraceletTicketCategoryRepository) FindByEventIDAndCategoryID(eventID string, categoryID string) (*domain.FindEventBraceletCategoryWithCategory, error) {
	logger := xlogger.Logger
	var eventBraceletTicketWithCategory domain.FindEventBraceletCategoryWithCategory

	query := `
    SELECT 
		ebc.id,
        btc.name AS category_name, 
        btc.is_used_several_time, 
        ebc.max_use_per_event
    FROM event_bracelet_categories ebc
    JOIN bracelet_ticket_categories btc ON ebc.category_id = btc.id
    WHERE ebc.event_id = ? AND ebc.category_id = ?
    LIMIT 1
`

	if err := m.db.Raw(query, eventID, categoryID).Scan(&eventBraceletTicketWithCategory).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to find event bracelet ticket category")
		return nil, err
	}

	return &eventBraceletTicketWithCategory, nil
}

// InsertEventBraceletTicketCategory implements domain.MysqlEventBraceletTicketCategoryRepository.
func (m *mysqlEventBraceletTicketCategoryRepository) InsertEventBraceletTicketCategory(eventBraceletTicketCategory domain.EventBraceletTicketCategory) error {
	panic("unimplemented")
}

// FindMaxUsePerEventByEventIDAndCategoryID implements domain.MysqlEventBraceletTicketCategoryRepository.
func (m *mysqlEventBraceletTicketCategoryRepository) FindMaxUsePerEventByEventIDAndCategoryID(ID string, eventID string) (int, error) {
	logger := xlogger.Logger
	var maxUsePerEvent int

	query := `
	SELECT 
		max_use_per_event
	FROM event_bracelet_categories
	WHERE id = ? AND event_id = ?
	LIMIT 1
`

	if err := m.db.Raw(query, ID, eventID).Scan(&maxUsePerEvent).Error; err != nil {
		logger.Error().Err(err).Msg("Failed to find max use per event")
		return 0, err
	}

	return maxUsePerEvent, nil
}
