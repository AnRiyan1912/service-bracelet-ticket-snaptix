package braceletticket

import (
	"encoding/json"
	"fmt"
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
	braceletTicketExelService            domain.BraceletTicketExelService
	cfg                                  *config.Config
}

func NewBraceletTicketService(braceletTicketRepository domain.MysqlBraceletTicketRepository, mysqlBraceletCategoryRepository domain.MysqlBraceletTicketCategoryRepository, mysqlTicketRepository domain.MysqlTicketRepository, mysqlEventRepository domain.MysqlEventRepository, mysqlEventBraceletCategoryRepository domain.MysqlEventBraceletTicketCategoryRepository, braceletTicketExelService domain.BraceletTicketExelService, cfg *config.Config) domain.BraceletTicketService {
	return &BraceletTicketService{
		mysqlBraceletTicketRepository:   braceletTicketRepository,
		mysqlBraceletCategoryRepository: mysqlBraceletCategoryRepository,
		mysqlTicketRepository:           mysqlTicketRepository, mysqlEventRepository: mysqlEventRepository,
		mysqlEventBraceletCategoryRepository: mysqlEventBraceletCategoryRepository,
		braceletTicketExelService:            braceletTicketExelService,
		cfg:                                  cfg,
	}
}

// CheckInBraceletTicket implements domain.BraceletTicketService.
func (b *BraceletTicketService) CheckInBraceletTicketOnline(eventId string, noTicketEncrypted string, deviceId string, deviceName string) (*domain.ApiResponseWithaoutData, error) {
	logger := xlogger.Logger

	// Find bracelet ticket
	getBraceletTicket, err := b.mysqlBraceletTicketRepository.FindByNoTicketEncrypted(noTicketEncrypted)
	if err != nil {
		if err.Error() == "record not found" {
			return &domain.ApiResponseWithaoutData{
				StatusCode: 404,
				Error:      true,
				Message:    "Bracelet ticket not found",
			}, nil
		}
		logger.Error().Err(err).Msg("Failed to find bracelet ticket")
		return nil, err
	}

	// Check if bracelet ticket is not found
	if getBraceletTicket == nil {
		logger.Error().Msg("Bracelet ticket not found")
		return &domain.ApiResponseWithaoutData{
			StatusCode: 404,
			Error:      true,
			Message:    "Bracelet ticket not found",
		}, nil
	}

	// Process decript bracelet ticket validate
	checkValid := utils.VerifyShortCode([]byte(getBraceletTicket.NoTicket), []byte(b.cfg.EncriptionKey), noTicketEncrypted)
	if !checkValid {
		logger.Error().Err(err).Msg("Failed to decrypt bracelet ticket")
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket not valid",
		}, nil
	}

	// check if bracelet ticket is valid same with event id
	if getBraceletTicket.EventID != eventId {
		logger.Error().Msg("Bracelet ticket not valid with event id")
		return &domain.ApiResponseWithaoutData{
			StatusCode: 403,
			Error:      true,
			Message:    "Bracelet ticket not valid",
		}, nil
	}

	// parse bracelet ticket session
	var braceletTicketSessions []domain.BraceletSession
	err = json.Unmarshal([]byte(getBraceletTicket.Sessions), &braceletTicketSessions)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to parse bracelet ticket session")
		return nil, fmt.Errorf("Failed to parse bracelet ticket session: %w", err)
	}

	// check session bracelet ticket
	responValidate, err := validateBracelet(braceletTicketSessions)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to validate bracelet ticket")
		return nil, err
	}
	if responValidate.Error {
		return responValidate, nil
	}
	// update bracelet ticket status
	err = b.mysqlBraceletTicketRepository.UpdateStatusDeviceIdAndNameById(getBraceletTicket.ID, deviceId, deviceName)
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

// CheckInBraceletTicketOffline implements domain.BraceletTicketService.
func (b *BraceletTicketService) CheckInBraceletTicketOffline(datas []domain.CheckInBraceletTicketOnlineRequest) error {
	logger := xlogger.Logger

	for _, data := range datas {
		// Check in bracelet ticket online
		// find bracelet ticket
		getBraceletTicket, err := b.mysqlBraceletTicketRepository.FindByNoTicketEncrypted(data.QrData)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to find bracelet ticket")
			return err
		}

		// update bracelet ticket status
		err = b.mysqlBraceletTicketRepository.UpdateStatusDeviceIdAndNameById(getBraceletTicket.ID, data.DeviceID, data.DeviceName)
		if err != nil {
			return err
		}

		// update total check in bracelet ticket
		err = b.mysqlEventRepository.UpdateTotalCheckInBraceletTicketByEventId(data.EventID, 1)
		if err != nil {
			return err
		}

	}

	return nil
}

// GenerateBraceletQrCode implements domain.BraceletTicketService.
func (b BraceletTicketService) GenerateBraceletQrCode(eventID string, braceletCategoryId string, total int, sessions []domain.BraceletSession) error {

	getEventBraceletCategory, err := b.mysqlEventBraceletCategoryRepository.FindByEventIDAndCategoryID(eventID, braceletCategoryId)
	if err != nil {
		return err
	}
	var braceletQrCodeDatas []domain.BraceletQrCodeData
	// prosess generate qr code and save bracelet ticket to database
	for i := 0; i < total; i++ {
		generateNoTicket := utils.GenerateRandomNoTicket(6)
		encriptNoTicket := utils.GenerateShortCode([]byte(generateNoTicket), []byte(b.cfg.EncriptionKey))
		if err != nil {
			return err
		}
		serialNumber := fmt.Sprintf("%0*d", len(fmt.Sprint(total)), i+1)
		braceletQrCodeData := domain.BraceletQrCodeData{
			No:       serialNumber,
			NoTicket: encriptNoTicket,
		}
		braceletQrCodeDatas = append(braceletQrCodeDatas, braceletQrCodeData)

		jsonSession, err := json.Marshal(sessions)
		if err != nil {
			return err
		}

		// generate serial number

		// insert bracelet ticket to database
		braceletTicket := domain.BraceletTicket{
			ID:                      utils.GenerateRandomId(),
			EventID:                 eventID,
			EventBraceletCategoryID: getEventBraceletCategory.ID,
			NoTicket:                generateNoTicket,
			Status:                  constan.NOT_YET_CHECK_IN,
			SerialNumber:            serialNumber,
			NoTicketEncrypted:       encriptNoTicket,
			Sessions:                string(jsonSession),
			CreatedAt:               utils.GetTimeNow(),
			UpdatedAt:               utils.GetTimeNow(),
		}

		err = b.mysqlBraceletTicketRepository.InsertBraceletTicket(braceletTicket)
		if err != nil {
			return err
		}
		// time.Sleep(1 * time.Second)
	}
	fileNameWithTime := fmt.Sprintf("%s%s", eventID, time.Now().Format("20060102150405"))
	// prosess generate excel file and save to database bracelet ticket exel
	err = b.braceletTicketExelService.GenerateBraceletTicketExel(braceletQrCodeDatas, eventID, fileNameWithTime)
	if err != nil {
		return err
	}
	// update total bracelet ticket in event
	err = b.mysqlEventRepository.UpdateTotalBraceletTicketByEventId(eventID, total)
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
