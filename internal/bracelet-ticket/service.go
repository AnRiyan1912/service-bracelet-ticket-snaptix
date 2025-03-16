package braceletticket

import (
	"encoding/json"
	"time"

	"bracelet-ticket-system-be/internal/constan"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/utils"
)

type BraceletTicketService struct {
	mysqlBraceletTicketRepository   domain.MysqlBraceletTicketRepository
	mysqlBraceletCategoryRepository domain.MysqlBraceletTicketCategoryRepository
	mysqlTicketRepository           domain.MysqlTicketRepository
	mysqlEventRepository            domain.MysqlEventRepository
	braceletQrCodeService           domain.BraceletQrCodeService
}

func NewBraceletTicketService(braceletTicketRepository domain.MysqlBraceletTicketRepository, mysqlBraceletCategoryRepository domain.MysqlBraceletTicketCategoryRepository, mysqlTicketRepository domain.MysqlTicketRepository, mysqlEventRepository domain.MysqlEventRepository, braceletQrCodeService domain.BraceletQrCodeService) domain.BraceletTicketService {
	return &BraceletTicketService{
		mysqlBraceletTicketRepository:   braceletTicketRepository,
		mysqlBraceletCategoryRepository: mysqlBraceletCategoryRepository,
		mysqlTicketRepository:           mysqlTicketRepository, mysqlEventRepository: mysqlEventRepository,
		braceletQrCodeService: braceletQrCodeService,
	}
}

// InsertBraceletTicket implements domain.BraceletTicketService.
func (b BraceletTicketService) InsertBraceletTicket(braceletTicket domain.BraceletTicket) error {
	panic("unimplemented")
}

// FindByBraceletTicketID implements domain.BraceletTicketService.
func (b BraceletTicketService) FindByBraceletTicketID(id int) (*domain.BraceletTicket, error) {
	panic("unimplemented")
}

// UpdateBraceletTicket implements domain.BraceletTicketService.
func (b BraceletTicketService) UpdateBraceletTicket(braceletTicket domain.BraceletTicket) error {
	panic("unimplemented")
}

// CheckInBraceletTicket implements domain.BraceletTicketService.
func (b BraceletTicketService) CheckInBraceletTicket(noTicket string) error {
	panic("unimplemented")
}

// GenerateBraceletQrCode implements domain.BraceletTicketService.
func (b BraceletTicketService) GenerateBraceletQrCode(eventID string, braceletCategoryId string, total int, sessions []domain.BraceletSession) error {

	braceletCategory, err := b.mysqlBraceletCategoryRepository.FindByBraceletTicketCategoryID(braceletCategoryId)
	if err != nil {
		return err
	}

	var braceletQrCodeDatas []domain.BraceletQrCodeData
	// prosess generate qr code and save bracelet ticket to database
	for range int(total) {
		braceletQrCodeData := domain.BraceletQrCodeData{
			EventID:  eventID,
			Sessions: sessions,
			NoTicket: utils.GenerateTicketNumber(8),
			Category: domain.BraceletTicketCategoryQrCodeData{
				Name:              braceletCategory.Name,
				IsUsedSeveralTime: braceletCategory.IsUsedSeveralTime,
				MaxUse:            braceletCategory.MaxUse,
			},
		}
		braceletQrCodeDatas = append(braceletQrCodeDatas, braceletQrCodeData)

		jsonSession, err := json.Marshal(sessions)
		if err != nil {
			return err
		}

		braceletTicket := domain.BraceletTicket{
			ID:                       utils.GenerateRandomId(),
			BraceletTicketCategoryID: braceletCategoryId,
			NoTicket:                 braceletQrCodeData.NoTicket,
			Status:                   constan.NOT_YET_CHECK_IN,
			Sessions:                 string(jsonSession),
			CreatedAt:                utils.GetTimeNow(),
			UpdatedAt:                utils.GetTimeNow(),
		}

		err = b.mysqlBraceletTicketRepository.InsertBraceletTicket(braceletTicket)
		if err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
	}
	folderName := eventID + "-" + utils.GetTimeNowWithoutSpacing()
	pathFolder, err := b.braceletQrCodeService.GenerateBraceletQrCode(braceletQrCodeDatas, folderName)
	if err != nil {
		return err
	}

	err = b.mysqlEventRepository.UpdateFolderNameTotalBraceletTicketByEventId(eventID, pathFolder, total)
	if err != nil {
		return err
	}
	return nil
}

// GetTotalBraceletAndTotalCheckInBraceletTicketByEventID implements domain.BraceletTicketService.
func (b *BraceletTicketService) GetTotalBraceletAndTotalCheckInBraceletTicketByEventID(eventID string) (*domain.GetTotalBraceletAndTotalCheckInBraceletTicketByEventIDRes, error) {
	getTotalBraceletAndTotalCheckInBraceletTicket, err := b.mysqlEventRepository.FindTotalAndTotalCheckInBraceletTicketByEventId(eventID)
	if err != nil {
		return nil, err
	}
	responTotal := domain.GetTotalBraceletAndTotalCheckInBraceletTicketByEventIDRes{
		TotalBraceletTicket: getTotalBraceletAndTotalCheckInBraceletTicket.TotalBraceletTicket,
		TotalCheckIn:        getTotalBraceletAndTotalCheckInBraceletTicket.TotalCheckInBraceletTicket,
	}
	return &responTotal, nil
}
