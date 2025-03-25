package braceletticket

import (
	"strconv"

	"gorm.io/gorm"

	"bracelet-ticket-system-be/internal/constan"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type MysqlBraceletTicketRepository struct {
	db *gorm.DB
}

func NewMysqlBraceletTicketRepository(db *gorm.DB) domain.MysqlBraceletTicketRepository {
	return MysqlBraceletTicketRepository{db: db}
}

// FindByBraceletTicketID implements domain.MysqlBraceletTicketRepository.
func (m MysqlBraceletTicketRepository) FindByNoTicketEncrypted(eventID string, noTicket string, eventBraceletCategoryID string) (*domain.BraceletTicket, error) {
	logger := xlogger.Logger
	var braceletTicket domain.BraceletTicket
	err := m.db.Where("event_id = ? AND no_ticket_encrypted = ? AND event_bracelet_category_id = ?", eventID, noTicket, eventBraceletCategoryID).First(&braceletTicket).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find bracelet ticket")
		return nil, err
	}
	return &braceletTicket, nil
}

// InsertBraceletTicket implements domain.MysqlBraceletTicketRepository.
func (m MysqlBraceletTicketRepository) InsertBraceletTicket(braceletTicket domain.BraceletTicket) error {
	logger := xlogger.Logger
	err := m.db.Create(&braceletTicket).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to insert bracelet ticket")
		return err
	}
	return nil
}

// UpdateBraceletTicket implements domain.MysqlBraceletTicketRepository.
func (m MysqlBraceletTicketRepository) UpdateStatusDeviceIdAndNameById(ID string, deviceID string, deviceName string) error {
	logger := xlogger.Logger
	err := m.db.Table("bracelet_tickets").Where("id = ?", ID).Updates(
		map[string]interface{}{
			"status":         constan.CHECKED_IN,
			"device_id":      deviceID,
			"device_name":    deviceName,
			"count_check_in": gorm.Expr("count_check_in + ?", 1),
		},
	).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update bracelet ticket")
	}
	return nil
}

// FindBySerialNumber implements domain.MysqlBraceletTicketRepository.
func (m MysqlBraceletTicketRepository) FindBySerialNumber(eventId string, serialNumber string, eventBraceletCategoryID string) (*domain.BraceletTicket, error) {
	logger := xlogger.Logger

	var braceletTicket domain.BraceletTicket
	err := m.db.Where("event_id = ? AND serial_number = ? AND event_bracelet_category_id = ?", eventId, serialNumber, eventBraceletCategoryID).First(&braceletTicket).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find bracelet ticket by serial number")
		return nil, err
	}
	return &braceletTicket, nil
}

// FindFirstWithLastSerialNumber implements domain.MysqlBraceletTicketRepository.
func (m MysqlBraceletTicketRepository) FindFirstWithLastSerialNumber(eventID string) (int, error) {
	logger := xlogger.Logger

	var braceletTicket domain.BraceletTicket
	err := m.db.Where("event_id = ?", eventID).Select("serial_number").Order("serial_number DESC").First(&braceletTicket).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find bracelet ticket by event id")
		return 0, err
	}

	parseToInt, err := strconv.Atoi(braceletTicket.SerialNumber)
	if err != nil {
		return 0, err
	}

	return parseToInt, nil
}
