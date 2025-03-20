package event

import (
	"gorm.io/gorm"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type mysqlEventRepository struct {
	db *gorm.DB
}

func NewMysqlEventRepository(db *gorm.DB) domain.MysqlEventRepository {
	return &mysqlEventRepository{db: db}
}

// UpdateFolderNameBraceletTicketByEventId implements domain.MysqlEventRepository.
func (m *mysqlEventRepository) UpdateTotalBraceletTicketByEventId(eventId string, additionalBraceletTicket int) error {
	logger := xlogger.Logger
	err := m.db.Table("events").Where("id = ?", eventId).Updates(map[string]interface{}{
		"total_bracelet_ticket": gorm.Expr("total_bracelet_ticket + ?", additionalBraceletTicket),
	}).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update folder name and total bracelet ticket by event id")
		return err
	}
	return nil
}

// UpdateTotalCheckInBraceletTicketByEventId implements domain.MysqlEventRepository.
func (m *mysqlEventRepository) UpdateTotalCheckInBraceletTicketByEventId(eventId string, totalCheckInBraceletTicket int) error {
	logger := xlogger.Logger
	err := m.db.Table("events").Where("id = ?", eventId).Update("total_check_in_bracelet_ticket", gorm.Expr("total_check_in_bracelet_ticket + ?", totalCheckInBraceletTicket)).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update total check in bracelet ticket by event id")
		return err
	}
	return nil
}

// FindTotalAndTotalCheckInBraceletTicketByEventId implements domain.MysqlEventRepository.
func (m *mysqlEventRepository) FindTotalAndTotalCheckInBraceletTicketByEventId(eventId string) (*domain.Event, error) {
	logger := xlogger.Logger
	var event domain.Event
	err := m.db.Table("events").Select("total_bracelet_ticket, total_check_in_bracelet_ticket").Where("id = ?", eventId).First(&event).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find total and total check in bracelet ticket by event id")
		return nil, err
	}
	return &event, nil
}
