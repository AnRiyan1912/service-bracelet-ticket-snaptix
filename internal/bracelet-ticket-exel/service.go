package braceletticketexel

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"

	"bracelet-ticket-system-be/internal/domain"
	"bracelet-ticket-system-be/internal/utils"
)

type BraceletTicketExelService struct {
	mysqlBraceletTicketExelRepository domain.MysqlBraceletTicketExelRepository
}

func NewBraceletTicketExelService(mysqlBraceletTicketExelRepository domain.MysqlBraceletTicketExelRepository) domain.BraceletTicketExelService {
	return &BraceletTicketExelService{
		mysqlBraceletTicketExelRepository: mysqlBraceletTicketExelRepository,
	}
}

// GenerateBraceletTicketExel implements domain.BraceletTicketExelService.
func (b *BraceletTicketExelService) GenerateBraceletTicketExel(datas []domain.BraceletQrCodeData, eventID string, fileName string) error {
	// Make new folder for the Excel file
	folderPath := filepath.Join("folder-bracelet-ticket-exel", eventID)
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed make the folder: %w", err)
	}

	// Buat file Excel baru
	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.SetSheetName("Sheet1", sheetName)

	// Header kolom
	headers := []string{"No", "Data"}
	for i, h := range headers {
		col := string(rune('A' + i)) // Perbaikan dari fmt.Sprint()
		f.SetCellValue(sheetName, col+"1", h)
	}

	// Fill data
	for i, data := range datas {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), data.No)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), data.NoTicket)
	}
	fileName = fmt.Sprintf("bracelet-ticket-%s.xlsx", fileName)
	filePath := filepath.Join(folderPath, fileName)

	// Save file Excel
	if err := f.SaveAs(filePath); err != nil {
		return fmt.Errorf("Failed save file Excel: %w", err)
	}

	// Save to database fileName by eventID
	braceletTicketExel := domain.BraceletTicketExel{
		ID:        uuid.New().String(),
		EventID:   eventID,
		FileName:  fileName,
		CreatedAt: utils.GetTimeNow(),
		UpdatedAt: utils.GetTimeNow(),
	}

	if err := b.mysqlBraceletTicketExelRepository.InsertBraceletTicketExel(braceletTicketExel); err != nil {
		return err
	}
	return nil
}

// GetBraceletTicketExelByEventID implements domain.BraceletTicketExelService.
func (b *BraceletTicketExelService) GetBraceletTicketExelByEventID(eventID string) ([]domain.BraceletTicketExel, error) {
	panic("unimplemented")
}

// DeleteBraceletTicketExelByEventID implements domain.BraceletTicketExelService.
func (b *BraceletTicketExelService) DeleteBraceletTicketExelByEventID(eventID string) error {
	panic("unimplemented")
}
