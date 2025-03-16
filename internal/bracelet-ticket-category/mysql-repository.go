package braceletticketcategory

import (
	"gorm.io/gorm"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type mysqlBraceletTicketCategoryRepository struct {
	db *gorm.DB
}

func NewMysqlBraceletTicketCategoryRepository(db *gorm.DB) domain.MysqlBraceletTicketCategoryRepository {
	return &mysqlBraceletTicketCategoryRepository{db: db}
}

// FindByBraceletTicketCategoryID implements domain.MysqlBraceletTicketCategoryRepository.
func (m *mysqlBraceletTicketCategoryRepository) FindByBraceletTicketCategoryID(id string) (*domain.BraceletTicketCategory, error) {
	logger := xlogger.Logger
	var braceletTicketCategory domain.BraceletTicketCategory
	err := m.db.Where("id = ?", id).First(&braceletTicketCategory).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find bracelet ticket category")
		return nil, err
	}
	return &braceletTicketCategory, nil
}

// InsertBraceletTicketCategory implements domain.MysqlBraceletTicketCategoryRepository.
func (m *mysqlBraceletTicketCategoryRepository) InsertBraceletTicketCategory(braceletTicketCategory domain.BraceletTicketCategory) error {
	logger := xlogger.Logger
	err := m.db.Create(&braceletTicketCategory).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to insert bracelet ticket category")
		return err
	}
	return nil
}

// UpdateBraceletTicketCategory implements domain.MysqlBraceletTicketCategoryRepository.
func (m *mysqlBraceletTicketCategoryRepository) UpdateBraceletTicketCategory(braceletTicketCategory domain.BraceletTicketCategory) error {
	logger := xlogger.Logger
	err := m.db.Save(&braceletTicketCategory).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update bracelet ticket category")
		return err
	}
	return nil
}
