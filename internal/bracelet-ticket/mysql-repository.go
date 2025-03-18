package braceletticket

import (
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
func (m MysqlBraceletTicketRepository) FindByNoTicket(noTicket string) (*domain.BraceletTicket, error) {
	logger := xlogger.Logger
	var braceletTicket domain.BraceletTicket
	err := m.db.Where("no_ticket = ?", noTicket).First(&braceletTicket).Error
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
func (m MysqlBraceletTicketRepository) UpdateStatusById(ID string) error {
	logger := xlogger.Logger
	err := m.db.Table("bracelet_tickets").Where("id = ?", ID).Update("status", constan.CHECKED_IN).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update bracelet ticket")
	}
	return nil
}
