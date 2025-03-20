package braceletticketexel

import (
	"gorm.io/gorm"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type MysqlBraceletTicketExel struct {
	db *gorm.DB
}

func NewMysqlBraceletTicketExelRepository(db *gorm.DB) domain.MysqlBraceletTicketExelRepository {
	return MysqlBraceletTicketExel{db: db}
}

// InsertBraceletTicketExel implements domain.MysqlBraceletTicketExelRepository.
func (m MysqlBraceletTicketExel) InsertBraceletTicketExel(braceletTicketExel domain.BraceletTicketExel) error {
	logger := xlogger.Logger
	err := m.db.Create(&braceletTicketExel).Error
	if err != nil {
		logger.Err(err).Msg("Failed to insert bracelet ticket exel")
		return err
	}
	return nil
}

// FindByEventID implements domain.MysqlBraceletTicketExelRepository.
func (m MysqlBraceletTicketExel) FindByEventID(eventID string) ([]domain.BraceletTicketExel, error) {
	logger := xlogger.Logger
	var braceletTicketExel []domain.BraceletTicketExel
	err := m.db.Where("event_id = ?", eventID).Find(&braceletTicketExel).Error
	if err != nil {
		logger.Err(err).Msg("Failed to find bracelet ticket exel by event id")
		return nil, err
	}
	return braceletTicketExel, nil
}

// DeleteByEventID implements domain.MysqlBraceletTicketExelRepository.
func (m MysqlBraceletTicketExel) DeleteByEventID(eventID string) error {
	logger := xlogger.Logger
	err := m.db.Where("event_id = ?", eventID).Delete(&domain.BraceletTicketExel{}).Error
	if err != nil {
		logger.Err(err).Msg("Failed to delete bracelet ticket exel by event id")
		return err
	}
	return nil
}
