package ticket

import (
	"gorm.io/gorm"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type mysqlTicketRepository struct {
	db *gorm.DB
}

func NewMysqlTicketRepository(db *gorm.DB) domain.MysqlTicketRepository {
	return &mysqlTicketRepository{db: db}
}

// FindEventSessionsByTicketIdID implements domain.MysqlTicketRepository.
func (m *mysqlTicketRepository) FindEventSessionsByTicketIdID(ticketID string) (*domain.Ticket, error) {
	logger := xlogger.Logger

	var ticket domain.Ticket
	err := m.db.Where("id = ?", ticketID).First(&ticket).Error
	if err != nil {
		logger.Error().Err(err).Msg("Failed to find ticket")
		return nil, err
	}
	return &ticket, nil
}
