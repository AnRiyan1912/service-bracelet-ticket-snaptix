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
func (m *mysqlEventRepository) UpdateFolderNameTotalBraceletTicketByEventId(eventId string, folderName string, totalBraceletTicket int) error {
	logger := xlogger.Logger
	err := m.db.Table("events").Where("id = ?", eventId).Updates(map[string]interface{}{
		"folder_name_ bracelet_ticket": folderName,
		"total_bracelet_ticket":        totalBraceletTicket,
	}).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update folder name bracelet ticket by event id")
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
