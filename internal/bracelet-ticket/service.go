package braceletticket

import (
	"encoding/json"
	"time"

	"bracelet-ticket-system-be/internal/config"
	"bracelet-ticket-system-be/internal/constan"
	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/utils"
	"bracelet-ticket-system-be/pkg/xlogger"
)

type BraceletTicketService struct {
	mysqlBraceletTicketRepository        domain.MysqlBraceletTicketRepository
	mysqlBraceletCategoryRepository      domain.MysqlBraceletTicketCategoryRepository
	mysqlTicketRepository                domain.MysqlTicketRepository
	mysqlEventRepository                 domain.MysqlEventRepository
	mysqlEventBraceletCategoryRepository domain.MysqlEventBraceletTicketCategoryRepository
	braceletQrCodeService                domain.BraceletQrCodeService
	cfg                                  *config.Config
}

func NewBraceletTicketService(braceletTicketRepository domain.MysqlBraceletTicketRepository, mysqlBraceletCategoryRepository domain.MysqlBraceletTicketCategoryRepository, mysqlTicketRepository domain.MysqlTicketRepository, mysqlEventRepository domain.MysqlEventRepository, mysqlEventBraceletCategoryRepository domain.MysqlEventBraceletTicketCategoryRepository, braceletQrCodeService domain.BraceletQrCodeService, cfg *config.Config) domain.BraceletTicketService {
	return &BraceletTicketService{
		mysqlBraceletTicketRepository:   braceletTicketRepository,
		mysqlBraceletCategoryRepository: mysqlBraceletCategoryRepository,
		mysqlTicketRepository:           mysqlTicketRepository, mysqlEventRepository: mysqlEventRepository,
		mysqlEventBraceletCategoryRepository: mysqlEventBraceletCategoryRepository,
		braceletQrCodeService:                braceletQrCodeService,
		cfg:                                  cfg,
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
func (b BraceletTicketService) CheckInBraceletTicket(eventId string, qrData string) (*domain.ApiResponseWithaoutData, error) {
	logger := xlogger.Logger
	decriptedData, err := utils.DecryptAESGCM(qrData, b.cfg.EncriptionKey)
	if err != nil {
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket not valid",
		}, nil
	}
	var braceletQrCodeData domain.BraceletQrCodeData
	err = json.Unmarshal([]byte(decriptedData), &braceletQrCodeData)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal bracelet qr code data")
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket not valid",
		}, nil
	}

	// check if bracelet ticket is valid same with event id
	if braceletQrCodeData.EventID != eventId {
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket not valid",
		}, nil
	}

	// check session bracelet ticket
	responValidate, err := validateBracelet(braceletQrCodeData.Sessions)
	if err != nil {
		return nil, err
	}
	if responValidate.Error {
		return responValidate, nil
	}

	// check bracelet ticket in database
	getBraceletTicket, err := b.mysqlBraceletTicketRepository.FindByNoTicket(braceletQrCodeData.NoTicket)
	if err != nil {
		return nil, err
	}
	if getBraceletTicket == nil {
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket not valid",
		}, nil
	} else if getBraceletTicket.Status == constan.CHECKED_IN {
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket already check in",
		}, nil
	}

	// update bracelet ticket status
	err = b.mysqlBraceletTicketRepository.UpdateStatusById(getBraceletTicket.ID)
	if err != nil {
		return nil, err
	}

	// update total check in bracelet ticket
	err = b.mysqlEventRepository.UpdateTotalCheckInBraceletTicketByEventId(eventId, 1)
	if err != nil {
		return nil, err
	}

	return &domain.ApiResponseWithaoutData{
		StatusCode: 200,
		Error:      false,
		Message:    "Success check in",
	}, nil
}

// GenerateBraceletQrCode implements domain.BraceletTicketService.
func (b BraceletTicketService) GenerateBraceletQrCode(eventID string, braceletCategoryId string, total int, sessions []domain.BraceletSession) error {

	getEventBraceletCategory, err := b.mysqlEventBraceletCategoryRepository.FindByEventIDAndCategoryID(eventID, braceletCategoryId)
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
			Category: domain.FindEventBraceletCategoryWithCategory{
				IsUsedSeveralTime: getEventBraceletCategory.IsUsedSeveralTime,
				MaxUsePerEvent:    getEventBraceletCategory.MaxUsePerEvent,
				CategoryName:      getEventBraceletCategory.CategoryName,
			},
		}
		braceletQrCodeDatas = append(braceletQrCodeDatas, braceletQrCodeData)

		jsonSession, err := json.Marshal(sessions)
		if err != nil {
			return err
		}

		braceletTicket := domain.BraceletTicket{
			ID:                      utils.GenerateRandomId(),
			EventBraceletCategoryID: getEventBraceletCategory.ID,
			NoTicket:                braceletQrCodeData.NoTicket,
			Status:                  constan.NOT_YET_CHECK_IN,
			Sessions:                string(jsonSession),
			CreatedAt:               utils.GetTimeNow(),
			UpdatedAt:               utils.GetTimeNow(),
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

func validateBracelet(sessions []domain.BraceletSession) (*domain.ApiResponseWithaoutData, error) {
	getTimeNow := time.Now()
	currentDate := getTimeNow.Format("2006-01-02")
	foundValidSession := false

	for _, session := range sessions {
		startTime, err := time.Parse("2006-01-02 15:04:05", session.StartTime)
		if err != nil {
			return &domain.ApiResponseWithaoutData{StatusCode: 400, Error: true, Message: "Invalid start time format"}, nil
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", session.EndTime)
		if err != nil {
			return &domain.ApiResponseWithaoutData{StatusCode: 400, Error: true, Message: "Invalid end time format"}, nil
		}

		// Check whather the session is valid for today
		if startTime.Format("2006-01-02") == currentDate {
			foundValidSession = true

			// Check whether the session is valid for the current time
			if getTimeNow.After(endTime) {
				return &domain.ApiResponseWithaoutData{
					StatusCode: 403,
					Error:      true,
					Message:    "Bracelet ticket not valid for this session",
				}, nil
			}
		}
	}

	// If no valid session is found
	if !foundValidSession {
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "No valid session for today",
		}, nil
	}

	// If the bracelet ticket is valid
	return &domain.ApiResponseWithaoutData{
		StatusCode: 200,
		Error:      false,
		Message:    "Bracelet ticket valid",
	}, nil
}
